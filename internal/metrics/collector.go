package metrics

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vtapaskar/brahma/internal/config"
	"github.com/vtapaskar/brahma/internal/splunk"
	"github.com/vtapaskar/brahma/internal/storage"
	"go.uber.org/zap"
)

type DeviceMetric struct {
	DeviceID   string                 `json:"device_id"`
	DeviceType string                 `json:"device_type"`
	Hostname   string                 `json:"hostname"`
	Timestamp  time.Time              `json:"timestamp"`
	MetricType string                 `json:"metric_type"`
	Data       map[string]interface{} `json:"data"`
}

type CrashReport struct {
	ID         string    `json:"id"`
	DeviceID   string    `json:"device_id"`
	Hostname   string    `json:"hostname"`
	Timestamp  time.Time `json:"timestamp"`
	ReportType string    `json:"report_type"`
	Content    []byte    `json:"content"`
	Filename   string    `json:"filename"`
}

type Collector struct {
	config       config.MetricsConfig
	splunkClient *splunk.Client
	s3Client     *storage.S3Client
	logger       *zap.Logger
	metricBuffer []DeviceMetric
	bufferMu     sync.Mutex
	stopChan     chan struct{}
}

func NewCollector(cfg config.MetricsConfig, splunkClient *splunk.Client, s3Client *storage.S3Client, logger *zap.Logger) *Collector {
	c := &Collector{
		config:       cfg,
		splunkClient: splunkClient,
		s3Client:     s3Client,
		logger:       logger,
		metricBuffer: make([]DeviceMetric, 0, cfg.BufferSize),
		stopChan:     make(chan struct{}),
	}

	go c.flushLoop()

	return c
}

func (c *Collector) CollectMetric(metric DeviceMetric) error {
	c.bufferMu.Lock()
	defer c.bufferMu.Unlock()

	metric.Timestamp = time.Now()
	c.metricBuffer = append(c.metricBuffer, metric)

	if len(c.metricBuffer) >= c.config.BufferSize {
		return c.flush()
	}

	return nil
}

func (c *Collector) CollectCrashReport(report CrashReport) (string, error) {
	report.ID = uuid.New().String()
	report.Timestamp = time.Now()

	s3Key := c.s3Client.GenerateKey(report.DeviceID, report.ReportType, report.ID, report.Filename)

	if err := c.s3Client.Upload(s3Key, report.Content); err != nil {
		c.logger.Error("Failed to upload crash report to S3",
			zap.String("device_id", report.DeviceID),
			zap.String("report_id", report.ID),
			zap.Error(err),
		)
		return "", err
	}

	metadata := map[string]interface{}{
		"report_id":   report.ID,
		"device_id":   report.DeviceID,
		"hostname":    report.Hostname,
		"report_type": report.ReportType,
		"filename":    report.Filename,
		"s3_key":      s3Key,
		"timestamp":   report.Timestamp,
	}

	if err := c.splunkClient.SendEvent("crash_report", metadata); err != nil {
		c.logger.Warn("Failed to send crash report metadata to Splunk",
			zap.String("report_id", report.ID),
			zap.Error(err),
		)
	}

	c.logger.Info("Crash report stored",
		zap.String("report_id", report.ID),
		zap.String("device_id", report.DeviceID),
		zap.String("s3_key", s3Key),
	)

	return report.ID, nil
}

func (c *Collector) CollectBacktrace(report CrashReport) (string, error) {
	report.ReportType = "backtrace"
	return c.CollectCrashReport(report)
}

func (c *Collector) flushLoop() {
	ticker := time.NewTicker(time.Duration(c.config.FlushInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.bufferMu.Lock()
			if len(c.metricBuffer) > 0 {
				if err := c.flush(); err != nil {
					c.logger.Error("Failed to flush metrics", zap.Error(err))
				}
			}
			c.bufferMu.Unlock()
		case <-c.stopChan:
			return
		}
	}
}

func (c *Collector) flush() error {
	if len(c.metricBuffer) == 0 {
		return nil
	}

	for _, metric := range c.metricBuffer {
		data, err := json.Marshal(metric)
		if err != nil {
			c.logger.Error("Failed to marshal metric", zap.Error(err))
			continue
		}

		var eventData map[string]interface{}
		json.Unmarshal(data, &eventData)

		if err := c.splunkClient.SendEvent(metric.MetricType, eventData); err != nil {
			c.logger.Error("Failed to send metric to Splunk",
				zap.String("device_id", metric.DeviceID),
				zap.Error(err),
			)
		}
	}

	c.logger.Info("Flushed metrics to Splunk", zap.Int("count", len(c.metricBuffer)))
	c.metricBuffer = c.metricBuffer[:0]

	return nil
}

func (c *Collector) Stop() {
	close(c.stopChan)

	c.bufferMu.Lock()
	defer c.bufferMu.Unlock()
	c.flush()
}

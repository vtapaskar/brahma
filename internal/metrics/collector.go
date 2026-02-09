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

type BaseMetric struct {
	UID        string    `json:"uid"`
	Timestamp  time.Time `json:"timestamp"`
	MetricType string    `json:"metric_type"`
}

type LogReport struct {
	ID          string    `json:"id"`
	DeviceUID   string    `json:"device_uid"`
	Timestamp   time.Time `json:"timestamp"`
	LogType     string    `json:"log_type"`
	ProcessTag  string    `json:"process_tag"`
	Version     string    `json:"version"`
	Content     []byte    `json:"-"`
	Filename    string    `json:"filename"`
	S3Key       string    `json:"s3_key"`
}

type LogMetadata struct {
	LogID      string    `json:"log_id"`
	DeviceUID  string    `json:"device_uid"`
	LogType    string    `json:"log_type"`
	ProcessTag string    `json:"process_tag"`
	Version    string    `json:"version"`
	Filename   string    `json:"filename"`
	S3Key      string    `json:"s3_key"`
	Timestamp  time.Time `json:"timestamp"`
}

type Collector struct {
	config       config.MetricsConfig
	splunkClient *splunk.Client
	s3Client     *storage.S3Client
	logger       *zap.Logger
	metricBuffer []interface{}
	bufferMu     sync.Mutex
	stopChan     chan struct{}
}

func NewCollector(cfg config.MetricsConfig, splunkClient *splunk.Client, s3Client *storage.S3Client, logger *zap.Logger) *Collector {
	c := &Collector{
		config:       cfg,
		splunkClient: splunkClient,
		s3Client:     s3Client,
		logger:       logger,
		metricBuffer: make([]interface{}, 0, cfg.BufferSize),
		stopChan:     make(chan struct{}),
	}

	go c.flushLoop()

	return c
}

func (c *Collector) CollectCPUStats(stats *CPUStats) error {
	stats.Timestamp = time.Now()
	return c.bufferMetric("cpu_stats", stats.UID, stats)
}

func (c *Collector) CollectProcessStats(stats *ProcessStats) error {
	stats.Timestamp = time.Now()
	return c.bufferMetric("process_stats", stats.UID, stats)
}

func (c *Collector) CollectMgmtNetworkStats(stats *MgmtNetworkStats) error {
	stats.Timestamp = time.Now()
	return c.bufferMetric("mgmt_network_stats", stats.UID, stats)
}

func (c *Collector) CollectRouterBaseState(state *RouterBaseState) error {
	state.Timestamp = time.Now()
	return c.bufferMetric("router_base_state", state.UID, state)
}

func (c *Collector) bufferMetric(metricType string, uid string, data interface{}) error {
	c.bufferMu.Lock()
	defer c.bufferMu.Unlock()

	c.metricBuffer = append(c.metricBuffer, data)

	if len(c.metricBuffer) >= c.config.BufferSize {
		return c.flush()
	}

	return nil
}

func (c *Collector) CollectCrashReport(report *LogReport) (string, error) {
	report.ID = uuid.New().String()
	report.Timestamp = time.Now()
	report.LogType = "crash"

	s3Key := c.s3Client.GenerateLogKey(report.DeviceUID, report.ID, "crash")
	report.S3Key = s3Key

	if err := c.s3Client.Upload(s3Key, report.Content); err != nil {
		c.logger.Error("Failed to upload crash report to S3",
			zap.String("device_uid", report.DeviceUID),
			zap.String("log_id", report.ID),
			zap.Error(err),
		)
		return "", err
	}

	metadata := LogMetadata{
		LogID:      report.ID,
		DeviceUID:  report.DeviceUID,
		LogType:    report.LogType,
		ProcessTag: report.ProcessTag,
		Version:    report.Version,
		Filename:   report.Filename,
		S3Key:      s3Key,
		Timestamp:  report.Timestamp,
	}

	if err := c.sendLogMetadata(&metadata); err != nil {
		c.logger.Warn("Failed to send crash report metadata to Splunk",
			zap.String("log_id", report.ID),
			zap.Error(err),
		)
	}

	c.logger.Info("Crash report stored",
		zap.String("log_id", report.ID),
		zap.String("device_uid", report.DeviceUID),
		zap.String("s3_key", s3Key),
	)

	return report.ID, nil
}

func (c *Collector) CollectBacktrace(report *LogReport) (string, error) {
	report.ID = uuid.New().String()
	report.Timestamp = time.Now()
	report.LogType = "backtrace"

	s3Key := c.s3Client.GenerateLogKey(report.DeviceUID, report.ID, "backtrace")
	report.S3Key = s3Key

	if err := c.s3Client.Upload(s3Key, report.Content); err != nil {
		c.logger.Error("Failed to upload backtrace to S3",
			zap.String("device_uid", report.DeviceUID),
			zap.String("log_id", report.ID),
			zap.Error(err),
		)
		return "", err
	}

	metadata := LogMetadata{
		LogID:      report.ID,
		DeviceUID:  report.DeviceUID,
		LogType:    report.LogType,
		ProcessTag: report.ProcessTag,
		Version:    report.Version,
		Filename:   report.Filename,
		S3Key:      s3Key,
		Timestamp:  report.Timestamp,
	}

	if err := c.sendLogMetadata(&metadata); err != nil {
		c.logger.Warn("Failed to send backtrace metadata to Splunk",
			zap.String("log_id", report.ID),
			zap.Error(err),
		)
	}

	c.logger.Info("Backtrace stored",
		zap.String("log_id", report.ID),
		zap.String("device_uid", report.DeviceUID),
		zap.String("s3_key", s3Key),
	)

	return report.ID, nil
}

func (c *Collector) sendLogMetadata(metadata *LogMetadata) error {
	eventData := map[string]interface{}{
		"log_id":      metadata.LogID,
		"device_uid":  metadata.DeviceUID,
		"log_type":    metadata.LogType,
		"process_tag": metadata.ProcessTag,
		"version":     metadata.Version,
		"filename":    metadata.Filename,
		"s3_key":      metadata.S3Key,
		"timestamp":   metadata.Timestamp,
	}

	return c.splunkClient.SendEvent("log_metadata", eventData)
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

		metricType := c.getMetricType(metric)
		uid := c.getUID(eventData)

		if err := c.splunkClient.SendEvent(metricType, eventData); err != nil {
			c.logger.Error("Failed to send metric to Splunk",
				zap.String("uid", uid),
				zap.String("metric_type", metricType),
				zap.Error(err),
			)
		}
	}

	c.logger.Info("Flushed metrics to Splunk", zap.Int("count", len(c.metricBuffer)))
	c.metricBuffer = c.metricBuffer[:0]

	return nil
}

func (c *Collector) getMetricType(metric interface{}) string {
	switch metric.(type) {
	case *CPUStats:
		return "cpu_stats"
	case *ProcessStats:
		return "process_stats"
	case *MgmtNetworkStats:
		return "mgmt_network_stats"
	case *RouterBaseState:
		return "router_base_state"
	default:
		return "unknown"
	}
}

func (c *Collector) getUID(data map[string]interface{}) string {
	if uid, ok := data["uid"].(string); ok {
		return uid
	}
	return ""
}

func (c *Collector) Stop() {
	close(c.stopChan)

	c.bufferMu.Lock()
	defer c.bufferMu.Unlock()
	c.flush()
}

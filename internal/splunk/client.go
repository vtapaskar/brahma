package splunk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vtapaskar/brahma/internal/config"
	"go.uber.org/zap"
)

type Client struct {
	config     config.SplunkConfig
	httpClient *http.Client
	logger     *zap.Logger
}

type Event struct {
	Time       int64                  `json:"time"`
	Host       string                 `json:"host,omitempty"`
	Source     string                 `json:"source,omitempty"`
	SourceType string                 `json:"sourcetype,omitempty"`
	Index      string                 `json:"index,omitempty"`
	Event      map[string]interface{} `json:"event"`
}

func NewClient(cfg config.SplunkConfig, logger *zap.Logger) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !cfg.UseTLS,
		},
	}

	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
		logger: logger,
	}
}

func (c *Client) SendEvent(eventType string, data map[string]interface{}) error {
	data["event_type"] = eventType

	event := Event{
		Time:       time.Now().Unix(),
		Source:     c.config.Source,
		SourceType: c.config.SourceType,
		Index:      c.config.Index,
		Event:      data,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	scheme := "http"
	if c.config.UseTLS {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s:%d/services/collector/event", scheme, c.config.Host, c.config.Port)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Splunk "+c.config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("splunk returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) SendBatch(events []Event) error {
	var buffer bytes.Buffer

	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			c.logger.Warn("Failed to marshal event in batch", zap.Error(err))
			continue
		}
		buffer.Write(payload)
		buffer.WriteByte('\n')
	}

	scheme := "http"
	if c.config.UseTLS {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s:%d/services/collector/event", scheme, c.config.Host, c.config.Port)

	req, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return fmt.Errorf("failed to create batch request: %w", err)
	}

	req.Header.Set("Authorization", "Splunk "+c.config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send batch request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("splunk returned non-OK status for batch: %d", resp.StatusCode)
	}

	return nil
}

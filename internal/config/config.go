package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server  ServerConfig  `json:"server"`
	GRPC    GRPCConfig    `json:"grpc"`
	Splunk  SplunkConfig  `json:"splunk"`
	S3      S3Config      `json:"s3"`
	Metrics MetricsConfig `json:"metrics"`
}

type ServerConfig struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type GRPCConfig struct {
	Address   string `json:"address"`
	Port      int    `json:"port"`
	EnableTLS bool   `json:"enable_tls"`
	CertFile  string `json:"cert_file"`
	KeyFile   string `json:"key_file"`
}

type SplunkConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Token      string `json:"token"`
	Index      string `json:"index"`
	Source     string `json:"source"`
	SourceType string `json:"source_type"`
	UseTLS     bool   `json:"use_tls"`
}

type S3Config struct {
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
	Prefix          string `json:"prefix"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Endpoint        string `json:"endpoint"`
}

type MetricsConfig struct {
	BufferSize    int      `json:"buffer_size"`
	FlushInterval int      `json:"flush_interval_seconds"`
	DeviceTypes   []string `json:"device_types"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.GRPC.Port <= 0 || c.GRPC.Port > 65535 {
		return fmt.Errorf("invalid grpc port: %d", c.GRPC.Port)
	}

	if c.Splunk.Host == "" {
		return fmt.Errorf("splunk host is required")
	}

	if c.S3.Bucket == "" {
		return fmt.Errorf("s3 bucket is required")
	}

	if c.S3.Region == "" {
		return fmt.Errorf("s3 region is required")
	}

	return nil
}

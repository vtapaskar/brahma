# Brahma

A Go microservice for collecting metrics from SONiC network devices, storing them in Splunk, and managing crash reports/backtraces in S3.

## Features

- **Metrics Collection**: Listens for metrics from SONiC devices via HTTP API
- **Splunk Integration**: Buffers and forwards metrics to Splunk HEC (HTTP Event Collector)
- **S3 Storage**: Stores crash reports and backtraces with unique IDs in S3
- **Configurable**: JSON-based configuration for all components

## Project Structure

```
brahma/
├── cmd/brahma/          # Main application entry point
├── internal/
│   ├── config/          # Configuration loading and validation
│   ├── metrics/         # Metrics collection and processing
│   ├── models/          # SONiC device data models
│   ├── server/          # HTTP server and API handlers
│   ├── splunk/          # Splunk HEC client
│   └── storage/         # S3 storage client
├── config.example.json  # Example configuration file
├── Dockerfile           # Container build file
├── Makefile             # Build automation
└── go.mod               # Go module definition
```

## Quick Start

### Prerequisites

- Go 1.21+
- Access to Splunk HEC endpoint
- AWS S3 bucket (or compatible storage)

### Build

```bash
make build
```

### Configure

Copy the example configuration and update with your settings:

```bash
cp config.example.json config.json
```

Edit `config.json` with your Splunk and S3 credentials.

### Run

```bash
make run
```

Or directly:

```bash
./bin/brahma -config config.json
```

## API Endpoints

### Health Check
```
GET /health
```

### Submit Metrics
```
POST /api/v1/metrics
Content-Type: application/json

{
  "device_id": "switch-01",
  "device_type": "leaf",
  "hostname": "leaf-switch-01.dc1",
  "metric_type": "interface",
  "data": {
    "name": "Ethernet1",
    "rx_bytes": 1234567890,
    "tx_bytes": 9876543210
  }
}
```

### Upload Crash Report
```
POST /api/v1/crash-report
Content-Type: multipart/form-data

device_id: switch-01
hostname: leaf-switch-01.dc1
report_type: crash
file: @crash_dump.log
```

### Upload Backtrace
```
POST /api/v1/backtrace
Content-Type: multipart/form-data

device_id: switch-01
hostname: leaf-switch-01.dc1
file: @backtrace.txt
```

## Configuration

| Section | Field | Description |
|---------|-------|-------------|
| server.address | Listen address | Default: 0.0.0.0 |
| server.port | Listen port | Default: 8080 |
| splunk.host | Splunk HEC host | Required |
| splunk.port | Splunk HEC port | Default: 8088 |
| splunk.token | HEC authentication token | Required |
| splunk.index | Target Splunk index | Required |
| splunk.use_tls | Enable TLS | Default: true |
| s3.region | AWS region | Required |
| s3.bucket | S3 bucket name | Required |
| s3.prefix | Key prefix for objects | Optional |
| metrics.buffer_size | Metrics buffer before flush | Default: 100 |
| metrics.flush_interval_seconds | Flush interval | Default: 30 |

## Docker

Build and run with Docker:

```bash
make docker
docker run -p 8080:8080 -v $(pwd)/config.json:/app/config.json brahma:dev
```

## License

MIT

# MetricNexus Client
A simple CLI client to talk to a MetricNexus server.

## Usage
```bash
go run . config.yaml
```

### Create Metric
```bash
go run . config.yaml CREATE demo 'my demo metric'
```

### Update Metric
```bash
go run . config.yaml UPDATE demo 125
```

### Read Metric
```bash
go run . config.yaml READ demo # prints 125
```

## Config
```yaml
host: 0.0.0.0
port: 4096
key: UnsafeKeyNumber1
```
# MetricNexus Server
Starts a MetricNexus server using a given configuration.

## Usage
```bash
go run . config.yaml
```

## Config
```yaml
host: 0.0.0.0
port: 4096
state: 
key: 
cert: 
keys:
- UnsafeKeyNumber1
- UnsafeKeyNumber2
- UnsafeKeyNumber3
```

Leaving `state` empty lets the server store the state in the same directory as the config, replacing its file extension with `.state.yaml`. 
Leaving `key` and `cert` empty lets the server create a self-signed certificate automatically. 

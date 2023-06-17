# MetricNexus

MetricNexus is a library that provides a centralized source of metrics for Prometheus to scrape. It eliminates the need for manual configuration updates, firewall rule creation, and other repetitive tasks when new machines join an environment or cluster.
The library consists of a server and a client, both leveraging the `fiber` framework for efficient performance.

## Key Features
- **Centralized Metrics**: The server acts as a single source of metrics, enabling you to collect and analyze statistics for the entire cluster's lifetime, rather than focusing on individual members of the cluster. 
- **Stateful Server**: Each server is stateful and automatically saves its current metrics every minute. In case of server restarts, the metrics are preserved, ensuring seamless processing continuity.
- **API Key Authentication**: The server is secured with API key authentication. Only authorized clients with valid API keys can access and manipulate the metrics.
- **Automatic Self-Signed Certificate**: When the server is started without a key file and a certificate file (either empty strings or both files do not exist), the library automatically generates a self-signed certificate. 

## Use Case Examples
Here are a couple of use case examples highlighting the versatility of MetricNexus:

### Monitoring Spiders
In my [spiders](https://github.com/toxyl/spider) use case, I deploy numerous web spiders across multiple servers. Rather than focusing on individual spider activity, my interest lies in understanding their collective performance over time. With this solution, even if the servers restart due to various reasons, the metrics such as kills and time wasted can seamlessly resume without resetting each time a server reboots.

### Multi-Application Server Metrics
Another use case involves a server running multiple applications that communicate with a single instance of this library instead of exposing their own Prometheus endpoints. Each application can utilize its unique metric prefix, ensuring metrics don't overwrite each other and making it effortless to differentiate metrics across applications.

## Security Considerations
Please note that MetricNexus currently relies solely on API keys for authentication and operates as HTTP server without TLS. It is essential to implement additional security measures, such as network-level security, to protect sensitive data. Contributions implementing alternative protocols  are welcome.

## Server
```golang
server := metrics.NewServer("127.0.0.1", 3000, "/tmp/state.yaml")
server.AddAPIKey("Hello World")

// Either start with your own TLS certificate 
panic(server.Start("my.key", "my.cert"))

// Or let MetricNexus automatically create self-signed key and cert
panic(server.Start("", ""))
```
The server exposes Prometheus metrics at `/__metrics` and provides a CRUD REST API for manipulating metrics.

## Client
```golang
// Simple uptime metric
apiKey := "Hello World"
allowSelfSigned := true // set to false if you want to error on self-signed certificates
client := metrics.NewClient("127.0.0.1", 3000, apiKey, allowSelfSigned)
if err := client.Create("uptime", "metric server uptime"); err != nil {
    fmt.Printf("Failed to create uptime metric: %s\n", err.Error())
}

for {
    s := time.Since(t).Seconds()
    if err := client.Update("uptime", s); err != nil {
        fmt.Printf("Failed to set uptime metric: %s\n", err.Error())
    }
    time.Sleep(5 * time.Second)
}
```

The client has the following methods:
| Method | Returns | Description |
| --- | --- | --- |
| `Create(key, description string)` | `error` | Creates the metric if it doesn't exist. |
| `Update(key string, value interface{})` | `error` | Set the metric to the given value (casted to `float64`). |
| `CreateUpdate(key, description string, value interface{})` | `error` | First creates and then sets the metric. |
| `Read(key string)` | `(float64, error)` | Reads the metric. If an error occurs it will be returned as the second value. |
| `Increment(key string)` | `error` | Increments the metric. |
| `Decrement(key string)` | `error` | Decrements the metric. |
| `Add(key string, value interface{})` | `error` | Add the given value to the metric. |
| `Subtract(key string, value interface{})` | `error` | Subtracts the given value from the metric. |
| `Delete(key string)` | `error` | Unregisters the metric and removes it from the known metrics. **WARNING**: Creating the metric again, but with a different description, will cause a crash!  |

## API
If you need to control metrics from a non-Go application, you can utilize the REST API:

| Endpoint | Returns | OK Status | Description |
| --- | --- | --- | --- |
| `POST /:metric` | | 201 | Creates a new metric with the provided key and uses the request body as its description. |
| `GET /:metric` | float64 | 200 | Retrieves and returns the value of the specified metric. |
| `PUT /:metric` | | 204 | Updates the specified metric with the value from the request body. |
| `PUT /:metric/inc` | | 204 | Increments the specified metric. |
| `PUT /:metric/dec` | | 204 | Decrements the specified metric. |
| `PUT /:metric/add` | | 204 | Adds the value from the request body to the specified metric. |
| `PUT /:metric/sub` | | 204 | Subtracts the value from the request body from the specified metric. |
| `DELETE /:metric` | | 204 | **DANGER!** Unregisters the specified metric and removes it from the known metric list. Re-adding the metric with a different description will cause a crash! |

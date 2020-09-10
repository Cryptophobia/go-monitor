## go-monitor

go-monitor is a tool that collects metrics on external urls and exposes metrics for Prometheus to track.

go-monitor uses:
* Gorilla web toolkit [mux](http://www.gorillatoolkit.org/pkg/mux) for the routing.
* [Prometheus](https://prometheus.io/) time-series database for collecting metrics/monitoring.
* [prometheus/client_golang](https://github.com/prometheus/client_golang) for registering collections of custom metrics for the Prometheus metrics collector.
* [tcnksm/go-httpstat](https://github.com/tcnksm/go-httpstat) for precise timing of the go client http request.
* [prometheus-middleware](https://github.com/albertogviana/prometheus-middleware) for collecting metrics about API routes.
* [docker-compose](https://github.com/docker/compose) for building, developing, and testing locally.

### Configuration:

go-monitor needs a config.json file with this minimal structure:

```json
{
  "URLs": [
    {
      "name": "VMware",
      "url": "https://www.vmware.com"
    }
  ]
}
```

You can also specify the config file path with the `-config` flag. For example:

```
go run main.go -config a/custom/path/my_config.json
```

Kubernetes pod command:
```
command: go run main.go -config /go/src/app/config.json
```

### Metrics exported:

```
Name:      "sample_external_url_response_ms",
Subsystem: "go_monitor",
Help:      "Response time in ms for the external url",
Buckets:   []float64{50, 75, 100, 125, 150, 175, 200, 225, 250, 275, 300, 325, 350, 375, 400, 500, 1000, 2000}
Type:      Histogram

Name:      "sample_external_url_up",
Subsystem: "go_monitor",
Help:      "Monitor URL: 0 if the url is down, 1 if the url is up",
Type:      Gauge
```

### Run go-monitor locally:

```bash
make start TAG=latest ENV=dev
```

Builds the go-monitor and then runs it in the background - http://localhost:5000/metrics is the only route exposed for prometheus.
Root http://localhost:5000/ is also exposed with a status page and version information.

### Testing go-monitor locally:

```bash
go test ./tool/tests/...
```

### Tearing down go-monitor and prometheus:

```bash
make stop TAG=latest ENV=dev
```

### Images of Sample Metrics Prometheus UI Queries:

![go_monitor_ms_bucket_image](../master/docs/images/go_monitor_ms_bucket_image.png?raw=true)

![go_monitor_ms_bucket_le_200](../master/docs/images/go_monitor_ms_bucket_le_200.png?raw=true)

![go_monitor_sample_external_url_up](../master/docs/images/go_monitor_sample_external_url_up.png?raw=true)

![histogram_quantile](../master/docs/images/histogram_quantile.png?raw=true)

### Original Problem:

Overview

Create a solution (in either Golang or Python) designed to run on a Kubernetes Cluster to monitor internet urls and provide Prometheus metrics, once completed please upload your solution to your github.com account.

Requirements

    A service written in python or golang that queries 2 urls (https://httpstat.us/503 & https://httpstat.us/200)
    The service will check the external urls (https://httpstat.us/503 & https://httpstat.us/200 ) are up (based on http status code 200) and response time in milliseconds
    The service will run a simple http service that produces  metrics (on /metrics) and output a prometheus format when hitting the service /metrics url

        Expected response format:

            sample_external_url_up{url="https://httpstat.us/503 "}  = 0
            sample_external_url_response_ms{url="https://httpstat.us/503 "}  = [value]
            sample_external_url_up{url="https://httpstat.us/200 "}  = 1
            sample_external_url_response_ms{url="https://httpstat.us/200 "}  = [value]

    Looking for:

        Code in python or golang
        Dockerfile to build image
        Kubernetes Deployment Specification to deploy Image to Kubernetes Cluster
        Unit Tests
        Good readme providing instructions for use
        Points awarded for screen shots where the metrics are ingested by Prometheus Server

    Artifacts should be uploaded to github.com

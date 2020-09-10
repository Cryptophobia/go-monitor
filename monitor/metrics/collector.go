package metrics

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	config "monitor/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	httpstat "github.com/tcnksm/go-httpstat"
)

var validStatusCodes = map[int]bool{
	http.StatusCreated: true,
	http.StatusOK:      true,
}

// URL This struct handles status of each url
type URL struct {
	Name         string
	Status       []bool
	ResponseTime int
	URL          string
}

// Collector is our struct for the collector worker
type Collector struct {
	Config config.Config
	URLs   []*URL
}

// Creates the go_monitor_sample_external_url_up metric
var (
	sampleURLupGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "sample_external_url_up",
		Subsystem: "go_monitor",
		Help:      "Monitor URL: 0 if the url is down, 1 if the url is up",
	},
		[]string{
			// Which url are we monitoring?
			"url",
		},
	)
)

// Creates the go_monitor_sample_external_url_response_ms
var (
	sampleURLresponseMS = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "sample_external_url_response_ms",
		Subsystem: "go_monitor",
		Help:      "Response time in ms for the external url",
		Buckets:   []float64{50, 75, 100, 125, 150, 175, 200, 225, 250, 275, 300, 325, 350, 375, 400, 500, 1000, 2000},
		// Buckets here start at 50ms and increment by 25 until 400ms
	},
		[]string{
			// Which url are we monitoring?
			"url",
		},
	)
)

// Run function is where all the prometheus data logic and checks go in here with timeouts & statuses.
func (c *Collector) Run() {
	// Defining the metrics we wish to expose on init
	// Promauto automatically registers the metric if this function is ran in main()
	ticker := time.NewTicker(time.Second * time.Duration(c.Config.CheckTimer.Interval))

	for {
		select {
		case <-ticker.C:
			c.CheckAll()
		}
	}
}

// CheckAll checks all of the URLs and then updates our metrics
func (c *Collector) CheckAll() {
	client := http.Client{
		Timeout: time.Second * time.Duration(c.Config.CheckTimer.Timeout),
	}

	for _, url := range c.Config.URLs {
		resp, latency, err := fetchURL(&client, url.URL)

		// Convert latency nanonseconds to milliseconds
		latencyMS := float64(latency / time.Millisecond)

		status := true

		if err != nil {
			log.Printf("Error on request for go-monitor %s: %v", url.Name, err)
			status = false
		}

		if bool(!validStatusCodes[resp]) {
			status = false
		}

		// Change this to update the prometheus data
		if status == true {
			sampleURLupGauge.WithLabelValues(url.URL).Set(1)
			sampleURLresponseMS.WithLabelValues(url.URL).Observe(latencyMS)
		} else {
			sampleURLupGauge.WithLabelValues(url.URL).Set(0)
			sampleURLresponseMS.WithLabelValues(url.URL).Observe(000)
		}

	}
}

// fetchURL fetches an url and gets are response
func fetchURL(client *http.Client, URL string) (statusCode int, latency time.Duration, e error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return 0, 0, err
	}

	// Create go-httpstat powered context and pass it to http.Request
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	res, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("Cannot fetch url %s: %v", URL, err)
	}

	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	result.End(time.Now())

	// Show results
	log.Printf("Latency of request is: %+v \n", result)

	log.Printf("Made a http request to: %v, Response: %v", URL, res)

	totalResultLatency := result.NameLookup + result.Connect + result.Pretransfer + result.StartTransfer

	return res.StatusCode, totalResultLatency, nil
}

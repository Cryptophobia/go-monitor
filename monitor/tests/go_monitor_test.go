package tests

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const urlMetrics = "http://localhost:5000/metrics"

// TestGoMonitorGetMetrics tests the /metrics route exposed to Prometheus
func TestGoMonitorGetMetrics(t *testing.T) {
	req, err := http.NewRequest("GET", urlMetrics, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// Check if response is OK Status
	assert.Equal(t, "200 OK", resp.Status)
	assert.Equal(t, true, strings.Contains(string(body), "go_monitor_sample_external_url_response_ms_bucket"))
	assert.Equal(t, true, strings.Contains(string(body), "go_monitor_sample_external_url_response_ms_count"))
	assert.Equal(t, true, strings.Contains(string(body), "go_monitor_sample_external_url_up"))
}

package router

import (
	"fmt"
	"net/http"
)

const MONITOR_VERSION = "1.0"

type MonitorSvc struct {
	//connection *db.Connection
	version string
}

func NewMonitorSvc() *MonitorSvc {
	s := &MonitorSvc{
		version: MONITOR_VERSION,
	}
	return s
}

func (s *MonitorSvc) Root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "go-monitor \n"+
		"Version: %v \n"+
		"Checks if everything is healthy and reports for Prometheus! \n", s.version)
}

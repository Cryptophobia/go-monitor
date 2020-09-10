package main

import (
	"flag"
	"log"
	config "monitor/config"
	r "monitor/router"
	"net/http"

	check "monitor/metrics"
)

func main() {
	configFile := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()

	log.Printf(*configFile)

	conf, err := config.CreateConfigurationFromFile(*configFile)

	if err != nil {
		log.Fatalf("Could not load config.json file: %v", err)
	}

	collector := check.Collector{Config: conf}

	// Initialize the collector and run as go thread
	go collector.Run()

	NewMonitor := r.NewMonitorSvc()
	routes := r.InitRoutes(NewMonitor)
	router := r.NewMonitorRouter(routes)
	http.ListenAndServe(":5000", router)
}

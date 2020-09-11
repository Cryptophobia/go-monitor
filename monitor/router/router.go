package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	mdlw "github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

// NewMonitorRouter create a new router and then return mux.Router
func NewMonitorRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Create the go-http-metrics middleware for prometheus
	mdlw := mdlw.New(mdlw.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})
	h := std.Handler("", mdlw, router)

	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		//w.WriteHeader(http.StatusOK)
		h.ServeHTTP(w, r)
		fmt.Fprintln(w, "ok")
	})

	// Feed the router all of the information from the routes struct
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

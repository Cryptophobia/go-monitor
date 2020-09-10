package router

import "net/http"

// Route struct for routing
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is slice of routes struct
type Routes []Route

// InitRoutes initializes the routes and return them all
func InitRoutes(s *MonitorSvc) Routes {
	return Routes{
		Route{
			"Root",
			"GET",
			"/",
			s.Root,
		},
	}
}

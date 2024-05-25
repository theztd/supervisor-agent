package exporter

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	Port = ":8080"
)

type Metrics struct {
	Mutex        sync.Mutex
	StartTime    int
	SupTasks     []string
	Database     []string
	LowestUptime int
}

type Server struct {
	Port           string
	BaseAuthPath   string
	Metrics        *Metrics
	HealthInterval int
}

func (s *Server) Run() {
	log.Println("INFO [exporter]: Server listen on http://0.0.0.0" + s.Port)
	log.Println("     /metrics/    --- prometheus exporter")
	log.Println("     /_healthz/   --- health check endpoint returns 200 / 503, depends on minimal job uptime", s.HealthInterval)

	http.HandleFunc("/metrics/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("INFO [exporter]: Request from ", r.RemoteAddr, " to ", r.URL.Path)
		w.Header().Set("Content-Type", "text/plain")

		uptime := int(time.Now().UnixMilli()) - s.Metrics.StartTime

		fmt.Fprint(w, "# HELP supervisord_agent_uptime Exporter's uptime in miliseconds.\n")
		fmt.Fprint(w, "# TYPE supervisord_agent_uptime gauge\n")
		fmt.Fprintf(w, "supervisord_agent_uptime{} %d\n", uptime)

		// Lock the metrics to prevent concurrent accesss
		s.Metrics.Mutex.Lock()
		// Write the metrics to the response writer
		for _, l := range s.Metrics.SupTasks {
			fmt.Fprint(w, l)
		}
		for _, l := range s.Metrics.Database {
			fmt.Fprint(w, l)
		}

		// Unlock the metrics
		s.Metrics.Mutex.Unlock()
	})

	// health check endpoint returning 200 only if the supervisor jobs runs longer than X seconds
	http.HandleFunc("/_healthz/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		// log.Println("DEBUG [http]: Lowest uptime is ", s.Metrics.LowestUptime, " and treshold is ", s.HealthInterval)
		if s.Metrics.LowestUptime > s.HealthInterval {
			w.WriteHeader(200)
			fmt.Fprintf(w, "ERR: The lowest uptime is %d\n", s.Metrics.LowestUptime)
		} else {
			w.WriteHeader(503)
			fmt.Fprintf(w, "OK: The lowest uptime is %d\n", s.Metrics.LowestUptime)
		}
	})

	log.Fatalln(http.ListenAndServe(s.Port, nil))
}

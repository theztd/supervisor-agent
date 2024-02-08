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
	Mutex     sync.Mutex
	StartTime int
	SupTasks  []string
	Database  []string
}

type Server struct {
	Port         string
	BaseAuthPath string
	Metrics      *Metrics
}

func (s *Server) Run() {
	log.Println("INFO [exporter]: Server listen on http://0.0.0.0" + s.Port + "/metrics/ ...")

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

	log.Fatalln(http.ListenAndServe(s.Port, nil))
}

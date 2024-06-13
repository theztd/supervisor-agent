package main

import (
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"theztd/supervisor-agent/checks"
	"theztd/supervisor-agent/exporter"
	"time"

	_ "github.com/lib/pq"
)

var (
	dsn               string
	port              string
	supervisorUrl     string
	pgScript          string
	checkInterval     int
	minimalHelthyTime int
	info              bool
)

func GetEnv(name, def string) string {
	/*
		If ENV is not defined, return default
	*/
	if eVal := os.Getenv(name); eVal != "" {
		return eVal
	} else {
		return def
	}
}

func main() {
	flag.StringVar(&checks.MetricsDir, "metrics-dir", GetEnv("METRICS_DIR", "./metrics"), "Directory where metrics will be stored (METRICS_DIR).")
	flag.StringVar(&dsn, "pg-dsn", GetEnv("PG_DSN", ""), "PostgreSQL DSN (PG_DSN). Example: \"user=username dbname=mydb sslmode=disable\"")
	flag.StringVar(&pgScript, "pg-script", GetEnv("PG_SCRIPT", ""), "Script that be executed when PostgreSQL is not available (PG_SCRIPT). Example: ./path_to_restart_script.sh")
	flag.StringVar(&port, "port", GetEnv("PORT", ":8080"), "Exporter listening port (PORT).")
	flag.StringVar(&supervisorUrl, "supervisor-url", GetEnv("SUPERVISOR_URL", "http://127.0.0.1:9001/RPC2"), "RPC Supervisor interface URL (SUPERVISOR_URL).")
	flag.IntVar(&checkInterval, "check-interval", 30, "Interval between checks in seconds.")
	flag.IntVar(&minimalHelthyTime, "health-uptime", 30, "Minimal jobs uptime in seconds to set healthz endpoint to healthy state (return 200).")
	flag.BoolVar(&info, "version", false, "Print information about version and exits (5).")
	flag.Parse()

	if info {
		fmt.Printf("supervisor-agent (version: %s)\n", VERSION)
		fmt.Println("")
		flag.PrintDefaults()
		os.Exit(5)
	}

	log.Printf("INFO: Starting application (version: %s)...", VERSION)
	log.Println("DEBUG [INIT]:", dsn)

	server := exporter.Server{
		Port:           port,
		BaseAuthPath:   "",
		Metrics:        &exporter.Metrics{},
		HealthInterval: minimalHelthyTime,
	}
	server.Metrics.StartTime = int(time.Now().UnixMilli())

	go checks.GetSupervisordJobsUptime(server.Metrics, supervisorUrl, time.Duration(checkInterval))

	if dsn != "" && pgScript != "" {
		go checks.PgPing(server.Metrics, dsn, pgScript, time.Duration(checkInterval))
	}

	// go server.Debug()
	server.Run()
}

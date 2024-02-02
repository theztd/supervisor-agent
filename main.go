package main

import (
	"flag"
	"log"
	"os"
	"theztd/supervisor-agent/checks"
	"theztd/supervisor-agent/exporter"
	"time"

	_ "github.com/lib/pq"
)

var (
	dsn           string
	port          string
	supervisorUrl string
	pgScript      string
	checkInterval int
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
	log.Println("INFO: Starting application...")

	flag.StringVar(&checks.MetricsDir, "metrics-dir", GetEnv("METRICS_DIR", "./metrics"), "Directory where metrics will be stored (METRICS_DIR).")
	flag.StringVar(&dsn, "pg-dsn", GetEnv("PG_DSN", ""), "PostgreSQL DSN (PG_DSN). Example: \"user=username dbname=mydb sslmode=disable\"")
	flag.StringVar(&pgScript, "pg-script", GetEnv("PG_SCRIPT", ""), "Script that be executed when PostgreSQL is not available (PG_SCRIPT). Example: ./path_to_restart_script.sh")
	flag.StringVar(&port, "port", GetEnv("PORT", ":8080"), "Exporter listening port (PORT).")
	flag.StringVar(&supervisorUrl, "supervisor-url", GetEnv("SUPERVISOR_URL", "http://127.0.0.1:9001/RPC2"), "RPC Supervisor interface URL (SUPERVISOR_URL).")
	flag.IntVar(&checkInterval, "check-interval", 30, "Interval between checks in seconds.")
	flag.Parse()

	log.Println("DEBUG [INIT]:", dsn)

	server := exporter.Server{
		Port:         port,
		BaseAuthPath: "",
		Metrics:      &exporter.Metrics{},
	}
	server.Metrics.StartTime = int(time.Now().UnixMilli())

	go checks.GetSupervisordJobsUptime(server.Metrics, supervisorUrl, 5)

	if dsn != "" && pgScript != "" {
		go checks.PgPing(server.Metrics, dsn, pgScript, 5)
	}

	server.Run()
}

package main

import (
	"flag"
	"log"
	"os"
	"strconv"
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

func main() {
	log.Println("INFO: Starting application...")

	flag.StringVar(&checks.MetricsDir, "metrics-dir", "./metrics", "Directory where metrics will be stored (METRICS_DIR).")
	flag.StringVar(&dsn, "pg-dsn", "", "PostgreSQL DSN (PG_DSN). Example: \"user=username dbname=mydb sslmode=disable\"")
	flag.StringVar(&pgScript, "pg-script", "", "Script that be executed when PostgreSQL is not available (PG_SCRIPT).")
	flag.StringVar(&port, "port", ":8080", "Exporter listening port (PORT).")
	flag.StringVar(&supervisorUrl, "supervisor-url", "http://127.0.0.1:9001/RPC2", "RPC Supervisor interface URL (SUPERVISOR_URL).")
	flag.IntVar(&checkInterval, "check-interval", 30, "Interval between checks in seconds (CHECK_INTERVAL).")
	flag.Parse()

	// Check if env variables are set
	if ev := os.Getenv("METRICS_DIR"); ev != "" {
		checks.MetricsDir = ev
	}
	if ev := os.Getenv("PG_DSN"); ev != "" {
		dsn = ev
	}
	if ev := os.Getenv("PORT"); ev != "" {
		port = ev
	}
	if ev := os.Getenv("SUPERVISOR_URL"); ev != "" {
		supervisorUrl = ev
	}
	if ev := os.Getenv("PG_SCRIPT"); ev != "" {
		pgScript = ev
	}
	if ev := os.Getenv("CHECK_INTERVAL"); ev != "" {
		checkInterval, _ = strconv.Atoi(ev)
	}

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

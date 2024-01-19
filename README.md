# supervisor-agent

Application exports supervisord managed job's uptime for prometheus. 

## Usage

```bash
Usage supervisor-agent:
  -check-interval int
        Interval between checks in seconds (CHECK_INTERVAL). (default 30)
  -metrics-dir string
        Directory where metrics will be stored (METRICS_DIR). (default "./metrics")
  -pg-dsn string
        PostgreSQL DSN (PG_DSN). (default "user=username dbname=mydb sslmode=disable")
  -pg-script string
        Script that be executed when PostgreSQL is not available (PG_SCRIPT). (default "./reload-jobs.sh")
  -port string
        Exporter listening port (PORT). (default ":8080")
  -supervisor-url string
        RPC Supervisor interface URL (SUPERVISOR_URL). (default "http://127.0.0.1:9001/RPC2")
```


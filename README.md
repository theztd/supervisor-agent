# supervisor-agent

Application exports supervisord managed job's uptime and a few basic metrics (RSS, CPU%, IO R/W) at **http://0.0.0.0:${PORT}/metrics/** url. The second function of this application is to monitor postgresql and run defined script (for example reload all services, ...) when postgresql is not available after 3 checks in row.

## Routes

```bash
/metrics/    ... prometheus exporter
/_healthz/   ... health check returning 200 / 503
```

## Usage

```bash
Usage supervisor-agent:
  -check-interval int
        Interval between checks in seconds. (default 30)
  -health-uptime int
        Minimal jobs uptime in seconds to set healthz endpoint to healthy state (return 200). (default 30)
  -metrics-dir string
        Directory where metrics will be stored (METRICS_DIR). (default "./metrics")
  -pg-dsn string
        PostgreSQL DSN (PG_DSN). Example: "user=username dbname=mydb sslmode=disable"
  -pg-script string
        Script that be executed when PostgreSQL is not available (PG_SCRIPT). Example: ./path_to_restart_script.sh
  -port string
        Exporter listening port (PORT). (default ":8080")
  -supervisor-url string
        RPC Supervisor interface URL (SUPERVISOR_URL). (default "http://127.0.0.1:9001/RPC2")
  -version
        Print information about version and exits (5).
```

 * The arguments -pg-dsn and -pg-script have to be defined to enable pgPing check.
 * The check interval is not possible to set via ENV variable.


## Example

Depoy this binary together with the supervisord (it is possible to run this script as a systemd service or as a supervisor job), write your ./reload-jobs.sh script doing all required stuff and configure your prometheus to gather the metrics...

### Script reloading jobs using database (./reload-jobs.sh)


This script will be run after the database wasn't been available 3 times in row. Example is included in the repository as [./reload-jobs.sh](./reload-jobs.sh)


### Create supervisord job

Example how to configure supervisor job running supervisor-agent with custom DSN definition and script path (you can also set PG_DSN env variable to reach the same result)

```toml
[program:SupervisorExporter]
command=/usr/local/bin/supervisor-agent --check-interval 15 --pg-dsn "user=develop password=developPassword dbname=develop sslmode=disable" -pg-script /usr/local/bin/reload-jobs.sh
directory=%(here)s
autorestart=True
user=supervisor
process_name=%(process_num)02d
numprocs=1
```

### Configure prometheus to gather metrics

```yaml
scrape_config:
- job_name: supervisor_agent
  honor_timestamps: true
  scrape_interval: 30s
  scrape_timeout: 10s
  metrics_path: /metrics
  scheme: http
  follow_redirects: true
  enable_http2: true
  static_configs:
  - targets:
    - 10.1.11.6:8080
    labels:
      env: prod
```
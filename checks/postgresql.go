package checks

import (
	"database/sql"
	"log"
	"os"
	"os/exec"
	"theztd/supervisor-agent/exporter"
	"time"
)

func PgPing(metrics *exporter.Metrics, dsn, shellCommand string, checkInterval time.Duration) {
	log.Println("INFO [checks.postgres]: Starting PgPing check...")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Println("ERR [checks.postgres]: Unable to reach database...")
	}
	defer db.Close()

	for {
		// At least 1 of 3 connection attempts shuld be available
		var pingState error

		results := []string{}
		results = append(results, "# HELP supervisord_agent_postgresql_ping Ping of PostgreSQL database 1 = UP\n")
		results = append(results, "# TYPE supervisord_agent_postgresql_ping Ping of PostgreSQL gauge\n")

		for i := 0; i < 3; i++ {
			pingState = db.Ping()
			if pingState == nil {
				results = append(results, "supervisord_agent_postgresql_ping{} 1\n")
				metrics.Mutex.Lock()
				metrics.Database = results
				metrics.Mutex.Unlock()
				break
			}
			time.Sleep(1 * time.Second)
		}

		if pingState != nil {
			log.Println("WARNING [checks.postgres]: Database is not available, running defined shell command.")
			results = append(results, "supervisord_agent_postgresql_ping{} 0\n")
			metrics.Mutex.Lock()
			metrics.Database = results
			metrics.Mutex.Unlock()

			cmd := exec.Command("sh", "-c", shellCommand)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				log.Println("ERR [checks.postgres]: Error during shell command execution.")
				log.Println("ERR [checks.postgres]: Running ", shellCommand, " with result: ", err)
			}
		}

		time.Sleep(checkInterval * time.Second)
	}
}

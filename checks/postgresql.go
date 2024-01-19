package checks

import (
	"database/sql"
	"log"
	"os"
	"os/exec"
	"time"
)

func PgPing(dsn, shellCommand string, checkInterval time.Duration) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Println("ERR [checks.postgres]: Unable to reach database...")
	}
	defer db.Close()
	log.Println("INFO [checks.postgres]: Connection to database established.")

	for {
		// At least 1 of 3 connection attempts shuld be available
		var pingState error

		metricsFile, err := os.Create(MetricsDir + "/postgresql.txt")
		metricsFile.WriteString("# HELP supervisord_agent_postgresql_ping Ping of PostgreSQL database 1 = UP\n")
		metricsFile.WriteString("# TYPE supervisord_agent_postgresql_ping Ping of PostgreSQL gauge\n")
		for i := 0; i < 3; i++ {
			pingState = db.Ping()
			if err == nil {
				metricsFile.WriteString("supervisord_agent_postgresql_ping{} 1\n")
				break
			}
			time.Sleep(1 * time.Second)
		}

		if pingState != nil {
			log.Println("WARNING [checks.postgres]: Database is not available, running defined shell command.")
			metricsFile.WriteString("supervisord_agent_postgresql_ping{} 0\n")
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

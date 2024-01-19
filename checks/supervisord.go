package checks

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/abrander/go-supervisord"
)

func GetSupervisordJobsUptime(url string, interval time.Duration) {
	log.Println("INFO [checks.supervisord]: Starting supervisord checks...")
	for {
		c, err := supervisord.NewClient(url)
		if err != nil {
			log.Println("ERR [checks.supervisord]: Unable to connect supervisord...")
			// panic(err.Error())
		}
		svcs, err := c.GetAllProcessInfo()
		if err != nil {
			log.Println("ERR [checks.supervisord]: Unable to get status of supervisord jobs...")
			log.Println(err.Error())
		}

		metricsFile, err := os.Create(MetricsDir + "/jobs.txt")
		if err != nil {
			log.Println("ERR [checks.supervisord]: Unable to open metrics file...")
			log.Println(err.Error())
		}

		metricsFile.WriteString("# HELP supervisord_agent_service_uptime Uptime of given service\n")
		metricsFile.WriteString("# TYPE supervisord_agent_service_uptime gauge\n")

		for _, svc := range svcs {
			svcUptime := svc.Now - svc.Start
			metricsFile.WriteString(fmt.Sprintf("supervisord_agent_service_uptime{service_name=\"%s:%s\", status=\"%s\" } %d\n", svc.Group, svc.Name, svc.StateName, svcUptime))
		}

		metricsFile.Close()
		time.Sleep(interval * time.Second)
	}
}

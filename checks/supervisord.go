package checks

import (
	"fmt"
	"log"
	"slices"
	"theztd/supervisor-agent/exporter"
	"time"

	"github.com/abrander/go-supervisord"
)

func GetSupervisordJobsUptime(metrics *exporter.Metrics, url string, interval time.Duration) {
	log.Println("INFO [checks.supervisord]: Starting supervisord checks...")

	for {
		uptimes := []int{}
		fmt.Println("UPTIMES:", uptimes)
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
		c.Close()

		results := []string{}
		results = append(results, "# HELP supervisord_agent_service_rss_bytes Memory consumption of process in Bytes\n")
		results = append(results, "# TYPE supervisord_agent_service_rss_bytes gauge\n")

		results = append(results, "# HELP supervisord_agent_service_cpu_percent CPU ussage in percent\n")
		results = append(results, "# TYPE supervisord_agent_service_cpu_percent gauge\n")

		results = append(results, "# HELP supervisord_agent_service_io_read_bytes IO Reads in bytes\n")
		results = append(results, "# TYPE supervisord_agent_service_io_read_bytes gauge\n")

		results = append(results, "# HELP supervisord_agent_service_io_read_count IO Read operations counter\n")
		results = append(results, "# TYPE supervisord_agent_service_io_read_count gauge\n")

		results = append(results, "# HELP supervisord_agent_service_io_write_bytes IO Writes in byte\n")
		results = append(results, "# TYPE supervisord_agent_service_io_write_bytes gauge\n")

		results = append(results, "# HELP supervisord_agent_service_io_write_count IO Write operations counter\n")
		results = append(results, "# TYPE supervisord_agent_service_io_write_count gauge\n")

		results = append(results, "# HELP supervisord_agent_service_uptime Uptime of given service\n")
		results = append(results, "# TYPE supervisord_agent_service_uptime gauge\n")

		for _, svc := range svcs {
			procInfo, err := GetProcessInfo(svc.Pid)
			//memBytes, err := GetMemoryUsageBytes(svc.Pid)
			if err != nil {
				log.Println("ERR [checks.supervisord]: Unable to get metrics for supervisord job", svc.Group, svc.Name, "with PID", svc.Pid)
				log.Println(err.Error())
			}
			// RSS
			results = append(
				results,
				fmt.Sprintf("supervisord_agent_service_rss_bytes{service_name=\"%s:%s\", status=\"%s\" } %d\n", svc.Group, svc.Name, svc.StateName, procInfo.Memory.RSS),
			)
			// CPU
			results = append(
				results,
				fmt.Sprintf("supervisord_agent_service_cpu_percent{service_name=\"%s:%s\", status=\"%s\" } %f\n", svc.Group, svc.Name, svc.StateName, procInfo.CPUPercent),
			)
			// IO read Bytes
			results = append(
				results,
				fmt.Sprintf("supervisord_agent_service_io_read_bytes{service_name=\"%s:%s\", status=\"%s\" } %d\n", svc.Group, svc.Name, svc.StateName, procInfo.IO.ReadBytes),
			)
			// IO read Count
			results = append(
				results,
				fmt.Sprintf("supervisord_agent_service_io_read_count{service_name=\"%s:%s\", status=\"%s\" } %d\n", svc.Group, svc.Name, svc.StateName, procInfo.IO.ReadCount),
			)
			// IO write Bytes
			results = append(
				results,
				fmt.Sprintf("supervisord_agent_service_io_write_bytes{service_name=\"%s:%s\", status=\"%s\" } %d\n", svc.Group, svc.Name, svc.StateName, procInfo.IO.WriteBytes),
			)
			// IO write Count
			results = append(
				results,
				fmt.Sprintf("supervisord_agent_service_io_write_count{service_name=\"%s:%s\", status=\"%s\" } %d\n", svc.Group, svc.Name, svc.StateName, procInfo.IO.WriteCount),
			)

			svcUptime := svc.Now - svc.Start
			results = append(
				results,
				fmt.Sprintf("supervisord_agent_service_uptime{service_name=\"%s:%s\", status=\"%s\" } %d\n", svc.Group, svc.Name, svc.StateName, svcUptime),
			)
			uptimes = append(uptimes, svcUptime)

			log.Println("INFO [checks.supervisord]: Svc", svc.Name, "uptime is", svcUptime)

		}
		results = append(results, "# HELP supervisord_agent_lowest_service_uptime The lowest uptime of all managed supervisord jobs\n")
		results = append(results, "# TYPE supervisord_agent_lowest_service_uptime gauge\n")
		results = append(
			results,
			fmt.Sprintf("supervisord_agent_lowest_uptime{} %d\n", slices.Min(uptimes)),
		)

		metrics.Mutex.Lock()
		metrics.SupTasks = results
		metrics.LowestUptime = slices.Min(uptimes)
		metrics.Mutex.Unlock()

		time.Sleep(interval * time.Second)
	}
}

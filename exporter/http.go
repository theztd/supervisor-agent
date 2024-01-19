package exporter

import (
	"log"
	"net/http"
	"os"
)

var (
	RootDir      = "/path/to/files"
	Port         = ":8080"
	baseAuthPath = ""
)

type Server struct {
	RootDir      string
	Port         string
	BaseAuthPath string
}

func (s *Server) Run() {
	log.Println("INFO [exporter]: Starting server...")
	log.Println("INFO [exporter]: Serving files from " + s.RootDir + " on port " + s.Port)

	http.Handle("/metrics/", basicAuth(http.StripPrefix("/metrics/", http.FileServer(http.Dir(s.RootDir)))))
	// http.Handle("/metrics/", basicAuth(http.FileServer(http.Dir(s.RootDir))))
	log.Fatalln(http.ListenAndServe(s.Port, nil))
}

func basicAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(baseAuthPath) < 1 {
			// Skip authentication
			handler.ServeHTTP(w, r)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok || !verifyUserAndPassword(username, password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized.\n"))
			log.Println("WARNING [exporter]: Unauthorized access to metrics.")
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func verifyUserAndPassword(user, password string) bool {
	// Read logins from htpasswd file
	file, err := os.Open(baseAuthPath)
	if err != nil {
		return false
	}
	defer file.Close()

	/*

		Implement ME

	*/

	return true
}

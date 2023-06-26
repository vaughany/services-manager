package main

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"os/exec"
	"regexp"

	"github.com/go-chi/chi/v5"
)

type config struct {
	services       []string
	serviceNames   map[string]string
	serviceResults map[string]bool
	regex          *regexp.Regexp
	webserver      struct {
		ipAddress string
		port      string
	}
}

func main() {
	var cfg config
	cfg.services = []string{"apache2", "mongod", "mysql", "postgresql"}
	cfg.serviceNames = map[string]string{
		"apache2":    "Apache",
		"mongod":     "MongoDB",
		"mysql":      "MySQL",
		"postgresql": "PostgreSQL",
	}
	cfg.serviceResults = make(map[string]bool)
	cfg.regex = regexp.MustCompile(`^active`)

	cfg.webserver.ipAddress = "localhost"
	cfg.webserver.port = "8888"

	log.Println("Services Manager starting")

	cfg.checkServices()

	router := chi.NewRouter()
	router.Get("/", cfg.indexHandler)
	router.Get("/stop/{service}", cfg.stopHandler)
	router.Get("/start/{service}", cfg.startHandler)
	router.Get("/all/{action}", cfg.allHandler)

	address := net.JoinHostPort(cfg.webserver.ipAddress, cfg.webserver.port)
	log.Printf("Web server running on http://%s\n", address)
	log.Fatal(http.ListenAndServe(address, router))
}

func (cfg *config) checkServices() {
	for _, service := range cfg.services {
		var out bytes.Buffer

		cmd := exec.Command("systemctl", "is-active", service)
		cmd.Stdout = &out
		cmd.Run()

		cfg.serviceResults[service] = cfg.regex.MatchString(out.String())
	}
}

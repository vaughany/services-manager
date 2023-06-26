package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"

	"github.com/go-chi/chi/v5"
)

func (cfg *config) indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("index handler")

	fmt.Fprint(w, "<h1>Services Manager</h1>\n")

	cfg.checkServices()

	fmt.Fprint(w, "<p><a href='/all/start'>Start all</a> or <a href='/all/stop'>stop all</a>.</p>\n")

	for _, service := range cfg.services {
		fmt.Fprintf(w, "<p>%s: ", cfg.serviceNames[service])
		if cfg.serviceResults[service] {
			fmt.Fprintf(w, "running (<a href='/stop/%s'>stop</a>)", service)
		} else {
			fmt.Fprintf(w, "stopped (<a href='/start/%s'>start</a>)", service)
		}
		fmt.Fprint(w, "</p>\n")
	}
}

func (cfg *config) stopHandler(w http.ResponseWriter, r *http.Request) {
	var (
		service = chi.URLParam(r, "service")
		out     bytes.Buffer
	)

	log.Printf("stop handler for %s\n", service)

	found := false
	for _, srv := range cfg.services {
		if service == srv {
			found = true
			break
		}
	}
	if !found {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		log.Printf("service %s not found\n", service)
		return
	}

	cmd := exec.Command("systemctl", "stop", service)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	log.Printf("%s stopped", service)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (cfg *config) startHandler(w http.ResponseWriter, r *http.Request) {
	var (
		service = chi.URLParam(r, "service")
		out     bytes.Buffer
	)

	log.Printf("start handler for %s\n", service)

	found := false
	for _, srv := range cfg.services {
		if service == srv {
			found = true
			break
		}
	}
	if !found {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		log.Printf("service %s not found\n", service)
		return
	}

	cmd := exec.Command("systemctl", "start", service)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	log.Printf("%s started", service)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (cfg *config) allHandler(w http.ResponseWriter, r *http.Request) {
	var (
		action = chi.URLParam(r, "action")
		wg     sync.WaitGroup
		out    bytes.Buffer
	)

	log.Printf("%s all handler\n", action)

	if action != "start" && action != "stop" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		log.Printf("action %s not found\n", action)
		return
	}

	for i, service := range cfg.services {
		wg.Add(1)

		go func(srv string, i int) {
			defer wg.Done()

			cmd := exec.Command("systemctl", action, srv)
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				panic(err)
			}

			log.Printf("> %d: %s %sed", i, srv, action)
		}(service, i)

	}

	wg.Wait()
	log.Println("all handler completed")

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

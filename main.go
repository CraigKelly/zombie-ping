package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

//TODO: working simple notifications
//TODO: if config missing, create on startup
//TODO: using config
//TODO: URL polling and status page
//TODO: actual notifications

var buildDate string // Set by our build script

func pcheck(err error) {
	if err != nil {
		log.Panicf("Fatal Error: %v\n", err)
	}
}

func jsonResponse(w http.ResponseWriter, req *http.Request, jsonSrc interface{}) {
	js, err := json.Marshal(jsonSrc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func runService(addrListen string) {
	// Read our config
	config, err := ReadConfiguration("./zombie-ping.json")
	pcheck(err)
	log.Printf("Read configuration: %d entries\n", len(config.Targets))

	// Start our URL checkers
	updateQuit := make(chan struct{})
	defer close(updateQuit)

	for _, t := range config.Targets {
		go func(target PingTarget) {
			seconds := time.Duration(target.PingSeconds) * time.Second
			updateTicker := time.NewTicker(seconds)
			checkNow := true
			for {
				if checkNow {
					//TODO: actual URL check
					log.Printf("Time to check %s\n", target.URL)
					checkNow = false
				}
				select {
				case <-updateTicker.C:
					checkNow = true
				case <-updateQuit:
					updateTicker.Stop()
					return
				}
			}
		}(t)
	}

	// TODO: any rest needed?
	// http.HandleFunc("/accts", func(w http.ResponseWriter, req *http.Request) {
	// 	log.Printf("GET %s - returning list of len %d\n", req.URL.Path, len(accts))
	// 	jsonResponse(w, req, accts)
	// })

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// provide some functions to our templates
	funcMap := template.FuncMap{
		"Year": func() string { return time.Now().Format("2006") },
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		//TODO: move out after debugging
		templates := template.Must(template.New("ui").Funcs(funcMap).ParseFiles("static/index.html"))

		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	if addrListen == "" {
		log.Printf("No host specified: using default\n")
		addrListen = "127.0.0.1:8142"
	}
	log.Printf("Starting listen on %s\n", addrListen)
	http.ListenAndServe(addrListen, nil)

	log.Printf("Exiting\n")
}

/////////////////////////////////////////////////////////////////////////////
// Entry point

func main() {
	log.Printf("STARTING zombie-ping - built %s\n", buildDate)

	flags := flag.NewFlagSet("zombie-ping", flag.ExitOnError)
	hostBinding := flags.String("host", "", "How to listen for service")
	pcheck(flags.Parse(os.Args))

	runService(*hostBinding)
}

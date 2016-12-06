package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var buildDate string // Set by our build script

/////////////////////////////////////////////////////////////////////////////
// Entry point

func main() {
	log.Printf("STARTING zombie-ping - built %s\n", buildDate)

	flags := flag.NewFlagSet("zombie-ping", flag.ExitOnError)
	hostBinding := flags.String("host", "", "How to listen for service")
	pcheck(flags.Parse(os.Args))

	runService(*hostBinding)
}

/////////////////////////////////////////////////////////////////////////////
// Actual web service/site

func runService(addrListen string) {
	// Read our config
	config, err := ReadConfiguration("./zombie-ping.json")
	pcheck(err)
	log.Printf("Read configuration: %d entries\n", len(config.Targets))

	//TODO: this is so ugly - we need to split subscription and status update
	//      into other files and make channel based

	// Subscription list
	// Endpoint for getting notification registrations from the browser
	//TODO: need to debounce within a 5 second period so that we aren't notifying too much
	subscriptions := make(map[string]bool)
	subscriptionsMtx := sync.RWMutex{}
	doNotify := func() {
		subscriptionsMtx.RLock()
		defer subscriptionsMtx.RUnlock()
		for u := range subscriptions {
			go func(url string) {
				log.Printf("Sending notification for %v\n", url)
				client := &http.Client{}
				req, _ := http.NewRequest("POST", url, nil)
				req.Header.Set("TTL", "60")
				client.Do(req)
			}(u)
		}
	}

	// Status map with safe update
	urlStatus := CreateURLStatusMap()
	urlUpdateMtx := sync.RWMutex{}
	doUpdate := func(url string, descrip string, good bool) {
		urlUpdateMtx.Lock()
		defer urlUpdateMtx.Unlock()
		urlStatus.SetState(url, descrip, good)
		if !good {
			doNotify()
		}
	}

	// Status checkers
	updateQuit := make(chan struct{})
	defer close(updateQuit)

	for _, t := range config.Targets {
		go func(target PingTarget) {
			seconds := time.Duration(target.PingSeconds) * time.Second
			updateTicker := time.NewTicker(seconds)
			checkNow := true
			for {
				if checkNow {
					log.Printf("Time to check %s\n", target.URL)
					resp, err := http.Get(target.URL)
					if err != nil {
						doUpdate(target.URL, fmt.Sprintf("Error: %v", err), false)
					} else {
						doUpdate(target.URL, resp.Status, resp.StatusCode == 200)
					}
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

	// Endpoint for getting current status list
	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		urlUpdateMtx.RLock()
		defer urlUpdateMtx.RUnlock()
		jsonResponse(w, req, urlStatus)
	})

	http.HandleFunc("/subscription-register", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		decoder := json.NewDecoder(req.Body)
		subscription := PushSubscription{}
		err := decoder.Decode(&subscription)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Got subscription: %v\n", subscription)
		subscriptionsMtx.Lock()
		defer subscriptionsMtx.Unlock()
		subscriptions[subscription.Endpoint] = true
	})

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// provide some functions to our templates
	funcMap := template.FuncMap{
		"Year": func() string { return time.Now().Format("2006") },
	}

	// create templates
	templates := template.Must(template.New("ui").Funcs(funcMap).ParseFiles("static/index.html"))

	// Template handlers
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Actually serve all our endpoints
	if addrListen == "" {
		log.Printf("No host specified: using default\n")
		addrListen = "127.0.0.1:8142"
	}
	log.Printf("Starting listen on %s\n", addrListen)
	http.ListenAndServe(addrListen, nil)

	log.Printf("Exiting\n")
}

/////////////////////////////////////////////////////////////////////////////
// Helper utilities

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

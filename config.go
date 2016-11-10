package main

import (
	"encoding/json"
	"io/ioutil"
)

// PingTarget is a single URL to check every PingSeconds
type PingTarget struct {
	URL         string
	PingSeconds int
}

// Configuration is the setup for our zombie pinger
type Configuration struct {
	Targets []PingTarget
}

// ReadConfiguration reads a config from a JSON file
func ReadConfiguration(filename string) (*Configuration, error) {
	contents, readerr := ioutil.ReadFile(filename)
	if readerr != nil {
		return nil, readerr
	}

	return DecodeConfiguration(contents)
}

// DecodeConfiguration reads JSON from a byte slice to construct a Configuration
func DecodeConfiguration(b []byte) (*Configuration, error) {
	config := &Configuration{}
	decodeerr := json.Unmarshal(b, config)
	if decodeerr != nil {
		return nil, decodeerr
	}
	return config, nil
}

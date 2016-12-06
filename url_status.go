package main

import (
	"time"
)

// URLStatus is the current status for a single URL
type URLStatus struct {
	StateDescrip string
	IsGood       bool
	LastCheck    string
}

// URLStatusMap a map of URL's to their status
type URLStatusMap map[string]*URLStatus

// CreateURLStatusMap creates an empty URL status map
func CreateURLStatusMap() URLStatusMap {
	return make(URLStatusMap)
}

// SetState update the status map for the given URL
func (m URLStatusMap) SetState(url string, state string, good bool) {
	lastCheck := time.Now().Format("Monday 3:04:05 PM - 1 Jan 2006")

	if urlState, present := m[url]; present {
		// Already had a record
		urlState.StateDescrip = state
		urlState.IsGood = good
		urlState.LastCheck = lastCheck
	} else {
		// New record
		m[url] = &URLStatus{
			StateDescrip: state,
			IsGood:       good,
			LastCheck:    lastCheck,
		}
	}
}

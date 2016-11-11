package main

// PushSubscription as defined at https://developer.mozilla.org/en-US/docs/Web/API/PushSubscription
type PushSubscription struct {
	Endpoint       string      `json:"endpoint"`
	SubscriptionID string      `json:"subscriptionId"`
	Options        interface{} `json:"options"`
}

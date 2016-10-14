package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MustWin/ctoll/ctoll/api/client"
	"github.com/MustWin/ctoll/ctoll/api/v1"
)

func main() {
	// ACME API key from ctoll/storage/seed/defaults.go
	const APIKEY = "2390511a-870d-11e6-ae22-56b6b6499611"
	const ENDPOINT = "http://localhost:9180"

	// Create a new client
	c := client.New(ENDPOINT, APIKEY, http.DefaultClient)
	fmt.Printf("created new client to %q with API key %q\n", ENDPOINT, APIKEY)

	// Check V1 endpoint is good and healthy
	//=====================================
	err := c.Ping()
	if err != nil {
		panic("error sending ping")
	}

	fmt.Println("sent ping")

	// Send a meter event
	//=====================================
	event := v1.StartMeterEvent{
		MeterEvent: &v1.MeterEvent{
			MeterID:   "meter-1",
			Timestamp: time.Now().Unix(),
			Type:      v1.MeterEventTypeStart,
			Container: &v1.ContainerInfo{
				ImageName: "tutum/hello-world",
				ImageTag:  "latest",
				Name:      "tutum-hello-world-example",
				Labels: map[string]string{
					"foo": "bar",
				},
			},
		},

		Allocated: &v1.BlockAlloc{
			VirtualCPUs: 0.5,
			MemoryMb:    4096.0,
		},
	}

	err = c.MeterEvents().SendStartMeter(event)
	if err != nil {
		panic("error sending event")
	}

	fmt.Println("sent meter event")
}

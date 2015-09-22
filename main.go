package main

import (
	"encoding/json"
	"github.com/samilton/peagent/checks"
	"log"
	"net/http"
)

func main() {
	alerts := make(chan checks.Message, 4096)
	log.Printf("Starting PE Agent")
	echo := checks.Random{Queue: alerts, Name: "Random Agent 007", Interval: 20}
	root := checks.Disk{Queue: alerts, Name: "Root Partition", Partition: "/", Interval: 5}
	tmp := checks.Disk{Queue: alerts, Name: "Tmp Partition", Partition: "/run", Interval: 15}

	go echo.Run()
	go root.Run()
	go tmp.Run()
	go CollectAlerts(alerts)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := checks.Message{}
		json.NewEncoder(w).Encode(m)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func CollectAlerts(messages chan checks.Message) {
	for {
		message := <-messages
		log.Printf("Alert Received: %T %t", message, message)
	}
}

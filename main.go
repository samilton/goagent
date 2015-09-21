package main

import (
	"encoding/json"
	"github.com/samilton/peagent/checks"
	"log"
	"net/http"
)

func main() {
	log.Printf("Starting PE Agent")
	echo := checks.Random{Name: "Random Agent 007", Interval: 20}
	root := checks.NewDiskCheck("Root Parition", "/", 10, .75)

	go echo.Run()
	go root.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := root.Report(root)
		json.NewEncoder(w).Encode(m)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

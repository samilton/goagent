package checks

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

var (
	currentStatus = CurrentStatus{"", ""}
	interval      = 1
)

type Random struct {
	Name     string
	Interval time.Duration
}

type CurrentStatus struct {
	Status string
	Value  string
}

func (e Random) Report() Message {
	log.Printf("Reporting %s status", e.Name)
	t := time.Now().Unix()
	m := Message{t, e.Name, currentStatus.Status, currentStatus.Value}

	return m
}

func (e Random) Run() {
	for {
		time.Sleep(e.Interval * time.Second)
		log.Printf("Performing %s Check", e.Name)
		currentStatus.Status = StatusClear
		currentStatus.Value = strconv.Itoa(rand.Intn(42))

	}
}

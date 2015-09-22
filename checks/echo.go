package checks

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Random struct {
	Queue    chan Message
	Name     string
	Interval time.Duration
}

func (e Random) Run() {
	for {
		time.Sleep(e.Interval * time.Second)
		log.Printf("Performing %s Check", e.Name)
		status := StatusClear
		value := strconv.Itoa(rand.Intn(42))

		t := time.Now().Unix()
		m := Message{t, e.Name, status, value}
		e.Queue <- m
	}
}

package checks

import (
	"github.com/samilton/peagent/types"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var (
	topic = "random"
)

type Random struct {
	Queue    chan types.Message
	Name     string
	Interval time.Duration
}

func (e Random) Run() {
	log.Printf("Performing %s Check", e.Name)
	for {
		time.Sleep(e.Interval * time.Second)
		log.Printf("Performing %s Check", e.Name)
		status := StatusClear
		value := strconv.Itoa(rand.Intn(42))

		t := time.Now().Unix()
		m := types.Message{t, topic, e.Name, status, value}

		e.Queue <- m
	}
}

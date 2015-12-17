package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"github.com/Shopify/sarama"
	"github.com/samilton/peagent/checks"
	"github.com/samilton/peagent/servers"
	"github.com/samilton/peagent/types"
	"log"
	"os"
	"strings"
	"time"
)

var (
	addr    = flag.String("addr", "8080", "The address to bind to")
	brokers = flag.String("brokers", "", "Comma-delimited list of brokers")
)

func main() {
	flag.Parse()
	log.Printf("Starting Peagent")

	brokerList := strings.Split(*brokers, ",")
	server := NewServer(brokerList)
	server.Run(*addr)
}

type Server struct {
	AlertQueue      chan types.Message
	ErrorQueue      chan error
	ChecksToPerform map[string]checks.Check
	AlertProducer   sarama.AsyncProducer
}

func ConfigureChecks(queue chan types.Message) (map[string]checks.Check, error) {

	checkConfig := make(map[string]checks.Check)
	checkConfig["echo"] = checks.Random{Queue: queue, Name: "test", Interval: 20}
	return checkConfig, nil
}

func NewServer(brokerList []string) *Server {
	queue := make(chan types.Message, 4096)
	producer := newAlertProducer(brokerList)
	checkConfiguration, err := ConfigureChecks(queue)
	if err != nil {
		panic(err)
	}
	server := &Server{
		AlertQueue:      queue,
		ChecksToPerform: checkConfiguration,
		AlertProducer:   producer,
	}

	return server
}

func (s *Server) Close() error {
	if err := s.AlertProducer.Close(); err != nil {
		log.Println("Failed to shut down alert producer cleanly", err)
	}

	return nil
}

func (s *Server) Run(addr string) error {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	for k, v := range s.ChecksToPerform {
		log.Printf("Starting %s check\n", k)
		go v.Run()
	}
	results := make(map[string]types.Message)

	go func() {
		for {
			select {
			case msg := <-s.AlertQueue:
				buf := new(bytes.Buffer)
				err := json.NewEncoder(buf).Encode(msg)
				if err != nil {
					s.ErrorQueue <- err
				}
				s.AlertProducer.Input() <- &sarama.ProducerMessage{
					Topic: msg.Topic,
					Key:   sarama.StringEncoder(hostname),
					Value: sarama.ByteEncoder(buf.Bytes()),
				}
			case err := <-s.ErrorQueue:
				panic(err)
			}
		}
	}()

	servers.HttpServer{Queue: results}.Run(addr)

	return nil
}

func createTlsConfiguration() (t *tls.Config) {
	return nil
}

func newAlertProducer(brokerList []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	tlsConfig := createTlsConfiguration()
	if tlsConfig != nil {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}

	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 + time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokerList, config)

	if err != nil {
		log.Fatalln("Failed to start Sarama Producer:", err)
	}

	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write alert:", err)
		}
	}()

	return producer
}

func (s Server) distributeAlerts(key string) {
}

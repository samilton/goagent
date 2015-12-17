package servers

import (
	"encoding/json"
	"github.com/samilton/peagent/types"
	"log"
	"net/http"
	"os"
	"time"
)

type Server interface {
	Run() string
}

type HttpServer struct {
	Queue map[string]types.Message
}

func (s HttpServer) Run(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := types.HostSummary{}
		t := time.Now().Unix()
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("Unable to determine hostname [%s]", err)
		} else {
			m.Hostname = hostname
		}
		m.Timestamp = t
		m.Messages = s.Queue
		json.NewEncoder(w).Encode(m)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

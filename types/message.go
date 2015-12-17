package types

type Message struct {
	// This is the Timestamp of the last update to this metric
	Timestamp int64 `json:"timestamp"`

	// The topic we will send on
	Topic string

	// The name of the check
	Name string `json:"name"`

	// Clear, Warn or Error
	Status string `json:"status"`

	// The value of the check
	Value string `json:"value"`
}

type Messages []Message

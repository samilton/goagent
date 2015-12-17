package types

type HostSummary struct {
	// This is the Timestamp of the last update to this metric
	Timestamp int64 `json:"timestamp"`

	Hostname string `json:"hostname"`

	Messages map[string]Message `json:"messages"`
}

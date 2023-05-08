package types

import (
	"time"
)

type OpenRPCDocument map[string]interface{}

// EventTopics represents topics for pub/sub events
type EventTopics string

const (
	// EventNodeOnline node online event
	EventNodeOnline EventTopics = "node_online"
	// EventNodeOffline node offline event
	EventNodeOffline EventTopics = "node_offline"
)

func (t EventTopics) String() string {
	return string(t)
}

// ValidationInfo Validation, election related information
type ValidationInfo struct {
	NextElectionTime time.Time
}

package chat_io_capi

import "encoding/json"

type Properties map[string]map[string]interface{}

// Events
type BaseEvent struct {
	ID         string     `json:"id"`
	Order      uint64     `json:"order"`
	Timestamp  uint64     `json:"timestamp,omitempty"`
	Type       string     `json:"type"`
	Properties Properties `json:"properties,omitempty"`
}

type EventMessage struct {
	*BaseEvent
	Text     string `json:"text"`
	AuthorID string `json:"author_id"`
	CustomID string `json:"custom_id,omitempty"`
}

// Messages
type IncomingEvent struct {
	ChatID   string          `json:"chat_id"`
	ThreadID string          `json:"thread_id"`
	Event    interface{}     `json:"-"`
	RawEvent json.RawMessage `json:"event"`
}

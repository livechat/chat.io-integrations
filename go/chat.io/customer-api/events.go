package chat_io_capi

import (
	"encoding/json"
	"errors"
)

type Events struct{}

func (*Events) Unmarshal(rawEvent json.RawMessage) (interface{}, error) {
	baseEvent := &BaseEvent{}
	if err := json.Unmarshal(rawEvent, baseEvent); err != nil {
		return nil, err
	}

	var event interface{}
	switch baseEvent.Type {
	case "message":
		event = &EventMessage{}
	default:
		return nil, errors.New("Unknown event type: " + baseEvent.Type)
	}

	if err := json.Unmarshal(rawEvent, event); err != nil {
		return nil, err
	}

	return event, nil
}

func (*Events) Type(e interface{}) string {
	base, ok := e.(*BaseEvent)
	if !ok {
		return ""
	}
	return base.Type
}

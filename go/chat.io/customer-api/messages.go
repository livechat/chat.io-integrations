package chat_io_capi

import "encoding/json"

type Messages struct{}

func (*Messages) IncomingEvent(raw json.RawMessage) *IncomingEvent {
	res := &IncomingEvent{}

	err := json.Unmarshal(raw, res)
	if err != nil {
		return nil
	}

	return res
}

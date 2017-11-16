package chat_io_capi

type Utils struct {
	events   *Events
	messages *Messages
}

func NewUtils() *Utils {
	return &Utils{
		events:   &Events{},
		messages: &Messages{},
	}
}

func (u *Utils) Events() *Events {
	return u.events
}

func (u *Utils) ProtocolMessages() *Messages {
	return u.messages
}

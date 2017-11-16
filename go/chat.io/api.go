package chat_io

import chat_io_capi "livechat/integration/go/chat.io/customer-api"

var (
	capi = &customerAPI{
		chat_io_capi.NewUtils(),
		chat_io_capi.NewRESTAPI(),
	}
)

type customerAPI struct {
	*chat_io_capi.Utils
	rest *chat_io_capi.RESTAPI
}

func (c *customerAPI) REST() *chat_io_capi.RESTAPI {
	return c.rest
}

func CustomerAPI() *customerAPI {
	return capi
}

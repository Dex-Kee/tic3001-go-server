package constant

const (
	// code
	RespCodeSuccess          = 200
	RespCodeResourceNotFound = 404
	RespCodeClientParamError = 400
	RespCodeServerError      = 500

	// message
	RespMsgSuccess          = "success"
	RespMsgResourceNotFound = "error: request resource is not found"
	RespMsgServerError      = "error: server error"
)

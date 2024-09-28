package service

const (
	ethAddrLength   = 42
	clusterIdLength = 64
	pubKeyLength    = 96
)
const (
	successCode           = 200
	badRequestCode        = 400
	notFoundCode          = 404
	serverErrCode         = 500
	badGatewayCode        = 502
	serverUnavailableCode = 503
)

const (
	successMsg           = "Success"
	badRequestMsg        = "Bad Request"
	notFoundMsg          = "Not Found"
	serverErrMsg         = "Internal Server Error"
	badGatewayMsg        = "Bad Gateway"
	serverUnavailableMsg = "Server Unavailable"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func newResponse(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}

var (
	successRes           = newResponse(successCode, successMsg)
	badRequestRes        = newResponse(badRequestCode, badRequestMsg)
	notFoundRes          = newResponse(notFoundCode, notFoundMsg)
	serverErrRes         = newResponse(serverErrCode, serverErrMsg)
	badGatewayRes        = newResponse(badGatewayCode, badGatewayMsg)
	serverUnavailableRes = newResponse(serverUnavailableCode, serverUnavailableMsg)
)

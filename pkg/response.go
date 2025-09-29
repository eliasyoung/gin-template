package pkg

const (
	CodeOk           = 0
	CodeError        = -1
	CodeUnauthorized = -2
)

type JSONResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type envelope struct {
	Msg   string `json:"msg"`
	MsgZh string `json:"msg_zh"`
}

func SuccessResponse(data interface{}) JSONResponse {
	return JSONResponse{
		Code: CodeOk,
		Data: data,
	}
}

func MessageResponse(code int, msg, msgZh string) JSONResponse {
	return JSONResponse{
		Code: code,
		Data: envelope{
			Msg:   msg,
			MsgZh: msgZh,
		},
	}
}

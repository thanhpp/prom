package dto

type Resp struct {
	Error RespError   `json:"error"`
	Data  interface{} `json:"data"`
}

func (r *Resp) SetCodeMsg(code int, msg string) {
	r.Error.Code = code
	r.Error.Message = msg
}

func (r *Resp) SetCode(code int) {
	r.Error.Code = code
}

func (r *Resp) SetMessage(msg string) {
	r.Error.Message = msg
}

func (r *Resp) SetData(data interface{}) {
	r.Data = data
}

type RespError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

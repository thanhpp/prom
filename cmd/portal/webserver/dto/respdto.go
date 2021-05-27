package dto

type Resp struct {
	Error RespError   `json:"error"`
	Data  interface{} `json:"data,omitempty"`
}

func (r *Resp) SetCodeMsg(code int, msg string) *Resp {
	r.Error.Code = code
	r.Error.Message = msg
	return r
}

func (r *Resp) SetCode(code int) *Resp {
	r.Error.Code = code
	return r
}

func (r *Resp) SetMessage(msg string) *Resp {
	r.Error.Message = msg
	return r
}

func (r *Resp) SetData(data interface{}) *Resp {
	r.Data = data
	return r
}

type RespError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

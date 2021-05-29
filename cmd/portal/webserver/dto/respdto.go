package dto

type Resp struct {
	RespError
	Data interface{} `json:"data,omitempty"`
}

func (r *Resp) SetData(data interface{}) *Resp {
	r.Data = data
	return r
}

type RespError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (r *RespError) SetCodeMsg(code int, msg string) *RespError {
	r.Error.Code = code
	r.Error.Message = msg
	return r
}

func (r *RespError) SetCode(code int) *RespError {
	r.Error.Code = code
	return r
}

func (r *RespError) SetMessage(msg string) *RespError {
	r.Error.Message = msg
	return r
}

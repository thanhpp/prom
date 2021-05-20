package errconst

import (
	"errors"
)

var (
	RPCEmptyRequestErr = errors.New("Empty RPC Request")
)

const (
	RPCEmptyReqCode = 1
	RPCDBErrorCode  = 2
	RPCSuccessCode  = 9
)

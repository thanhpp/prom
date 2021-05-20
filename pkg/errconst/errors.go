package errconst

import (
	"fmt"
)

type ServiceError struct {
	Srv  string `json:"service"`
	Err  error  `json:"error"`
	Code int32  `json:"code"`
	Msg  string `json:"message"`
}

func (e ServiceError) Error() string {
	if len(e.Msg) != 0 && e.Code != 0 {
		return fmt.Sprintf("ERROR - Service: %s. Error: %v. Code: %d. Message: %s",
			e.Srv, e.Err, e.Code, e.Msg)
	}

	return fmt.Sprintf("ERROR - Service: %s. Error: %v", e.Srv, e.Err)
}

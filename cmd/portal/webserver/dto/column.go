package dto

import (
	"github.com/thanhpp/prom/pkg/ccmanrpc"
)

type CreateNewColumnReq struct {
	ColumnName string `json:"columnName"`
}

type UpdateColumnIndex struct {
	ColumnID  uint32 `json:"columnID"`
	NextOfIdx uint32 `json:"nextOfIndex"`
}

type UpdateColumnReq struct {
	ColumnID uint32           `json:"columnID"`
	Column   *ccmanrpc.Column `json:"column"`
}

type DeleteColumn struct {
	ColumnID uint32 `json:"columnID"`
}

package dto

type CreateNewColumnReq struct {
	ColumnName string `json:"columnName"`
}

type UpdateColumnIndex struct {
	ColumnIndex []uint32 `json:"columnIndex"`
}

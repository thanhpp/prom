package dto

type CreateNewColumnReq struct {
	ColumnName string `json:"columnName"`
}

type UpdateColumnIndex struct {
	ColumnIndex []uint32 `json:"columnIndex"`
}

type DeleteColumn struct {
	ColumnID uint32 `json:"columnID"`
	MoveTo   uint32 `json:"moveToColumnID"`
}

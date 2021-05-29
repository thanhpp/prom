package dto

type CreateNewColumnReq struct {
	ColumnName string `json:"columnName"`
}

type UpdateColumnIndex struct {
	ColumnID  uint32 `json:"columnID"`
	NextOfIdx uint32 `json:"nextOfIndex"`
}

type DeleteColumn struct {
	ColumnID uint32 `json:"columnID"`
}

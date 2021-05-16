package dto

type CreateNewCardReq struct {
	ColumnID uint32 `json:"columnID"`
	Card     struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		AssignedTo  uint32 `json:"assignedTo"`
		DueDate     uint32 `json:"duedate"`
	} `json:"card"`
}

type UpdateCardInfoReq struct {
	ColumnID uint32 `json:"columnID"`
	Card     struct {
		ID          uint32 `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		AssignedTo  uint32 `json:"assignedTo"`
		DueDate     uint32 `json:"duedate"`
	} `json:"card"`
}

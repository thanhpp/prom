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

type ReorderCardOneColumnReq struct {
	ColumnID  uint32   `json:"columnID"`
	CardIndex []uint32 `json:"cardIndex"`
}

// ReorderCard if move in same col, colID = 0
type ReorderCard struct {
	CardID   uint32 `json:"cardID"`
	AboveIdx uint32 `json:"aboveOfIdx"`
	ColID    uint32 `json:"columnID"`
}

type DeleteCardReq struct {
	CardID uint32 `json:"cardID"`
}

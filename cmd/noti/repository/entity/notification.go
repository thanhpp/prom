package entity

import "time"

type Notification struct {
	ID        int        `json:"id" gorm:"Column:id; Type:int; PRIMARY KEY; AutoIncrement"`
	CardID    int        `json:"cardID" gorm:"Column:card_id; Type:int"`
	Content   string     `json:"content" gorm:"Column:content; Type:text; NOT NULL"`
	UserNotis []UserNoti `json:"-" gorm:"foreignKey:NotiID"`
	Seen      bool       `json:"seen" gorm:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt time.Time  `json:"deletedAt"`
}

type UserNoti struct {
	UserID    int       `json:"userID" gorm:"Column:user_id; Type:int; PRIMARY KEY"`
	NotiID    int       `json:"notificationID" gorm:"Column:noti_id; Type:int; PRIMARY KEY"`
	Seen      bool      `json:"seen" gorm:"Column:seen; Type:boolean; default:false"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

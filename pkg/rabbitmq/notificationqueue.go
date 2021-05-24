package rabbitmq

const (
	NotificationQueue = "notification"
)

type NewNotificationMsg struct {
	CardID  int    `json:"cardID"`
	Content string `json:"content"`
	UserIDs []int  `json:"userIDs"`
}

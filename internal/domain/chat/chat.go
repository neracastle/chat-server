package chat

import "time"

type Chat struct {
	Id        int64
	UserIds   []int64
	CreatedAt time.Time
}

func NewChat(userIds []int64) Chat {
	return Chat{
		UserIds:   userIds,
		CreatedAt: time.Now(),
	}
}

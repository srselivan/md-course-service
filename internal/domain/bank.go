package domain

import "time"

type Bank struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
	UserId         int64     `json:"userId"`
	QuestionsCount int64     `json:"questionsCount"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

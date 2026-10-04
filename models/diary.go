package models

import "time"

type DiaryEntries struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Date        time.Time `json:"date"`
	Subject     string    `json:"subject"`
	LessonTopic string    `json:"lesson_topic"`
	IsDone      bool      `json:"is_done"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Homework struct {
	ID      int    `json:"id"`
	DiaryID int    `json:"diary_id"`
	Subject string `json:"subject"`
}

package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title       string
	Description string
	Status      string    // backlog, todo, in_progress, done, archived
	AssigneeID  uuid.UUID // исполнитель (user_id)
	CreatorID   uuid.UUID // создатель задачи
	DueDate     *time.Time
	Labels      []string `gorm:"type:text[]"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

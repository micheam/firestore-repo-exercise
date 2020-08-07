package todolist

import (
	"time"

	"github.com/google/uuid"
)

//go:generate firestore-repo Task

// Meta ...
type Meta struct {
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	UpdatedBy string     `json:"-"`
	DeletedAt *time.Time `json:"-"`
	DeletedBy string     `json:"-"`
	Version   int        `json:"-"`
}

// Task ...
type Task struct {
	Meta
	ID   string `firestore:"-" firestore_key:"" json:"id"`
	Desc string `firestore:"description" json:"desc"`
	Done bool   `firestore:"done" json:"done"`
}

// NewTask create new (unsaved) task.
func NewTask(desc string) *Task {
	return &Task{
		ID:   uuid.New().String(),
		Desc: desc,
		Done: false,
	}
}

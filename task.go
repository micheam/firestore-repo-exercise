package todolist

import (
	"time"
)

//go:generate firestore-repo Task

// Meta ...
type Meta struct {
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt *time.Time
	DeletedBy string
	Version   int
}

// Task ...
type Task struct {
	Meta
	ID   string `firestore:"-" firestore_key:""`
	Desc string `firestore:"description"`
	Done bool   `firestore:"done"`
}

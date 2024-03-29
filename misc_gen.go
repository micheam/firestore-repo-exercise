// THIS FILE IS A GENERATED CODE. DO NOT EDIT
// generated version: 0.4.0
package todolist

import "time"

func SetLastThreeToZero(t time.Time) time.Time {
	return time.Unix(t.Unix(), int64(t.Nanosecond()/1000*1000))
}

type GetOption struct {
	IncludeSoftDeleted bool
}

type DeleteMode string

const (
	DeleteModeSoft = "soft"
	DeleteModeHard = "hard"
)

type DeleteOption struct {
	Mode DeleteMode
}

package structs

import (
	"time"
)

type Parameters struct {
	Filename string
	From     time.Time
	To       time.Time
}

type LogRecord struct {
	EventTime time.Time `json:"eventTime"`
	Email     string    `json:"email"`
	SessionId string    `json:"sessionId"`
}

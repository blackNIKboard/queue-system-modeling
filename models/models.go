package models

import "time"

type Request struct {
	IsFinished bool
	AppendTime time.Time
	EndTime    time.Time
}

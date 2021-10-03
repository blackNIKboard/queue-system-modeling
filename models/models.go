package models

import "time"

type Request struct {
	IsFinished bool
	AppendTime time.Duration
	EndTime    time.Duration
}

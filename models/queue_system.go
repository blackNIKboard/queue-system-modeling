package models

import "time"

type QSystem interface {
	Start() error
	Stop() error
	SendRequest(*Request) error
	GetAvgTime() time.Duration
	GetSystemTime() time.Duration
	CountQueuedRequests() int
	GetProcessedRequests() *[]Request
}

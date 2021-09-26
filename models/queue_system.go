package models

import "time"

type QSystem interface {
	Start(timeout int) error
	Stop() error
	SendRequest(request Request) error
	GetAvgTime() time.Duration
	CountQueuedRequests() int
	GetProcessedRequests() *[]Request
}

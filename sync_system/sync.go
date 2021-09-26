package sync_system

import (
	"context"
	"log"
	"time"

	"github.com/blackNIKboard/queue-system-modeling/models"
	gq "github.com/phf/go-queue/queue"
)

var _ models.QSystem = (*SyncSystem)(nil)

type SyncSystem struct {
	ctx     context.Context
	cancel  context.CancelFunc
	queue   *gq.Queue
	discard *[]models.Request
}

func NewSyncSystem() *SyncSystem {
	ctx, cancel := context.WithCancel(context.Background())

	return &SyncSystem{
		ctx:     ctx,
		cancel:  cancel,
		queue:   gq.New(),
		discard: &[]models.Request{},
	}
}

func (s SyncSystem) Start(timeout int) error {
	if timeout > 0 {
		s.ctx, s.cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	}

	go s.process()

	return nil
}

func (s SyncSystem) Stop() error {
	s.cancel()

	return nil
}

func (s SyncSystem) SendRequest(request models.Request) error {
	s.queue.PushBack(request)

	return nil
}

func (s SyncSystem) GetAvgTime() time.Duration {
	var sum time.Duration

	for _, request := range *s.discard {
		sum += request.EndTime.Sub(request.AppendTime)
	}

	return sum / time.Duration(len(*s.discard))
}

func (s SyncSystem) CountQueuedRequests() int {
	return s.queue.Len()
}

func (s SyncSystem) GetProcessedRequests() *[]models.Request {
	return s.discard
}

func (s SyncSystem) GetCtx() context.Context {
	return s.ctx
}

func (s SyncSystem) process() {
	for {
		select {
		case <-s.ctx.Done():
			log.Println("cancelled")

			return
		default:
			request := s.queue.PopFront().(models.Request)
			time.Sleep(time.Second)
			request.EndTime = time.Now()
			request.IsFinished = true

			*s.discard = append(*s.discard, request)
		}
	}
}

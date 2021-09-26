package async_system

import (
	"context"
	"log"
	"time"

	"github.com/blackNIKboard/queue-system-modeling/models"
	gq "github.com/phf/go-queue/queue"
)

var _ models.QSystem = (*AsyncSystem)(nil)

type AsyncSystem struct {
	ctx     context.Context
	cancel  context.CancelFunc
	queue   *gq.Queue
	discard *[]models.Request
}

func NewAsyncSystem() *AsyncSystem {
	ctx, cancel := context.WithCancel(context.Background())

	return &AsyncSystem{
		ctx:     ctx,
		cancel:  cancel,
		queue:   gq.New(),
		discard: &[]models.Request{},
	}
}

func (s AsyncSystem) Start(timeout int) error {
	if timeout > 0 {
		s.ctx, s.cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	}

	go s.process()

	return nil
}

func (s AsyncSystem) Stop() error {
	s.cancel()

	return nil
}

func (s AsyncSystem) SendRequest(request models.Request) error {
	s.queue.PushBack(request)

	return nil
}

func (s AsyncSystem) GetAvgTime() time.Duration {
	var sum time.Duration

	for _, request := range *s.discard {
		sum += request.EndTime.Sub(request.AppendTime)
	}

	return sum / time.Duration(len(*s.discard))
}

func (s AsyncSystem) CountQueuedRequests() int {
	return s.queue.Len()
}

func (s AsyncSystem) GetProcessedRequests() *[]models.Request {
	return s.discard
}

func (s AsyncSystem) GetCtx() context.Context {
	return s.ctx
}

func (s AsyncSystem) process() {
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

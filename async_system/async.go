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
	ctx         context.Context
	systemTime  *time.Duration
	stopIfEmpty bool
	cancel      context.CancelFunc
	queue       *gq.Queue
	discard     *[]models.Request
}

func NewAsyncSystem(timeout int, stopIfEmpty bool) *AsyncSystem {
	ctx, cancel := context.WithCancel(context.Background())

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	}

	duration := time.Duration(0)

	return &AsyncSystem{
		ctx:         ctx,
		cancel:      cancel,
		stopIfEmpty: stopIfEmpty,
		systemTime:  &duration,
		queue:       gq.New(),
		discard:     &[]models.Request{},
	}
}

func (s *AsyncSystem) Start() error {
	go s.process()

	return nil
}

func (s *AsyncSystem) Stop() error {
	s.cancel()

	return nil
}

func (s *AsyncSystem) SendRequest(request *models.Request) error {
	if request != nil {
		s.queue.PushBack(*request)

		return nil
	}

	s.queue.PushBack(models.Request{
		IsFinished: false,
		AppendTime: *s.systemTime,
		EndTime:    0,
	})

	return nil
}

func (s *AsyncSystem) GetAvgTime() time.Duration {
	var sum time.Duration

	for _, request := range *s.discard {
		sum += request.EndTime - request.AppendTime
	}

	return sum / time.Duration(len(*s.discard))
}

func (s *AsyncSystem) GetSystemTime() time.Duration {
	return *s.systemTime
}

func (s *AsyncSystem) CountQueuedRequests() int {
	return s.queue.Len()
}

func (s *AsyncSystem) GetProcessedRequests() *[]models.Request {
	return s.discard
}

func (s *AsyncSystem) GetCtx() context.Context {
	return s.ctx
}

func (s *AsyncSystem) process() {
	for {
		select {
		case <-s.ctx.Done():
			log.Println("cancelled")

			return
		default:
			if s.queue.Len() == 0 {
				if s.stopIfEmpty {
					s.cancel()
				}

				continue
			}

			request := s.queue.PopFront().(models.Request)
			//time.Sleep(time.Second)
			if request.AppendTime > *s.systemTime {
				s.queue.PushFront(request)
				*s.systemTime += time.Millisecond

				continue
			}

			*s.systemTime += time.Second

			request.EndTime = *s.systemTime
			request.IsFinished = true

			//fmt.Printf("processed, took e-a=r: %v-%v=%v\n", request.EndTime.Milliseconds(), request.AppendTime.Milliseconds(), (request.EndTime - request.AppendTime).Milliseconds())

			*s.discard = append(*s.discard, request)
		}
	}
}

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
	ctx         context.Context    // context of running modelling
	systemTime  *time.Duration     // modelling system time (relative (from 0 to ...))
	stopIfEmpty bool               //flag of stopping modelling if no requests in queue
	cancel      context.CancelFunc // cancelFunc of ctx
	systemQueue *gq.Queue          // overall requests queue
	queue       *gq.Queue          // 'real-time' queue
	discard     *[]models.Request  // processed requests

	userStat *Stat // statistics of users waiting for processing

	currentRequest *models.Request // current request being processed
}

type Stat struct {
	data  []int // not-zero values
	nulls int   //number of zeros occurred
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
		systemQueue: gq.New(),
		queue:       gq.New(),
		discard:     &[]models.Request{},
		userStat:    &Stat{},
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
		s.systemQueue.PushBack(*request)

		return nil
	}

	s.systemQueue.PushBack(models.Request{
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
	return s.queue.Len() + s.systemQueue.Len()
}

func (s *AsyncSystem) GetAvgUsers() float64 {
	var sum int

	for _, i := range s.userStat.data {
		sum += i
	}
	return float64(sum) / float64(len(s.userStat.data)+s.userStat.nulls)
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
			if s.systemQueue.Len() != 0 {
				newRequest := s.systemQueue.PopFront().(models.Request)
				if newRequest.AppendTime > *s.systemTime {
					s.systemQueue.PushFront(newRequest)
				} else {
					s.queue.PushBack(newRequest)
					//s.stat = append(s.stat, s.queue.Len())
				}
			}

			if s.queue.Len() == 0 {
				if s.stopIfEmpty && s.systemQueue.Len() == 0 {
					s.cancel()
				}
			} else if s.currentRequest == nil {
				request := s.queue.PopFront().(models.Request)

				s.currentRequest = &request
				s.currentRequest.EndTime = *s.systemTime + 500*time.Millisecond
			}

			if s.currentRequest != nil && (*s.currentRequest).EndTime < *s.systemTime {
				request := *s.currentRequest
				s.currentRequest = nil

				request.IsFinished = true
				//fmt.Printf("processed, took e-a=r: %v-%v=%v\n", request.EndTime.Milliseconds(), request.AppendTime.Milliseconds(), (request.EndTime - request.AppendTime).Milliseconds())
				*s.discard = append(*s.discard, request)

				if s.queue.Len() != 0 {
					s.userStat.data = append(s.userStat.data, s.queue.Len())
				} else {
					s.userStat.nulls++
				}
			}

			*s.systemTime += time.Millisecond
		}
	}
}

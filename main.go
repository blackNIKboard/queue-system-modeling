package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/blackNIKboard/queue-system-modeling/async_system"
	"github.com/blackNIKboard/queue-system-modeling/models"
	"github.com/davecgh/go-spew/spew"
)

func exp(lambda float64) float64 {
	return math.Log(1-rand.Float64()) / (-1 / lambda)
}

func avg(arr []float64) float64 {
	var sum float64

	for _, f := range arr {
		sum += f
	}

	return sum / float64(len(arr))
}

func waitForCtx(ctx context.Context, system *async_system.AsyncSystem) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Printf("waiting for the system: %10d processed, remaining %10d\n", len(*system.GetProcessedRequests()),
				system.CountQueuedRequests())
			time.Sleep(time.Second)
		}
	}
}

func testExp() {
	var arr []float64

	for i := 0; i < 10000; i++ {
		//log(1-rand.Float64())/(-λ)
		tmp := exp(1)
		//spew.Dump(tmp)
		arr = append(arr, tmp)
	}

	spew.Dump("---")
	spew.Dump(avg(arr))
}

func compute(alpha float64, requests int) (avgTime float64, avgUsers float64, overallTime float64) {
	log.Printf("---computing fo alpha %f\n", alpha)

	ss := async_system.NewAsyncSystem(50, true)

	// Pushing requests to system
	var prevRequest models.Request

	for i := 0; i < requests; i++ {
		var request models.Request

		if i == 0 {
			request.AppendTime = time.Duration(0) * time.Second
		} else {
			tmp := exp(1/alpha) * 1000
			request.AppendTime = prevRequest.AppendTime + time.Duration(tmp)*time.Millisecond
		}

		//fmt.Printf("Request #%d diff: %d\n", i, (request.AppendTime - prevRequest.AppendTime).Milliseconds())
		prevRequest = request

		ss.SendRequest(&request)
	}

	if err := ss.Start(); err != nil {
		return
	}

	waitForCtx(ss.GetCtx(), ss)

	avgTime = ss.GetAvgTime().Seconds()

	overallTime = ss.GetSystemTime().Seconds()
	avgUsers = ss.GetAvgUsers()

	log.Printf("avgTime %f, overallTime %f, avgUsers %f", avgTime, overallTime, avgUsers)

	return
}

func theorCompute(alpha float64) (avgTime float64, avgUsers float64) {
	avgUsers = alpha * (2 - alpha) / (2 * (1 - alpha))
	avgTime = avgUsers / alpha

	return
}

func main() {
	rand.Seed(time.Now().UnixNano())

	file, err := os.Create("res.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileTheor, err := os.Create("resTheor.txt")
	if err != nil {
		panic(err)
	}
	defer fileTheor.Close()

	for i := 0.01; i <= 1; i += 0.1 {
		avgTime, avgUsers, _ := compute(i, 100000)
		theorAvgTime, theorAvgUsers := theorCompute(i)
		fmt.Fprintf(file, "%5f %5f %5f\n", avgUsers, avgTime, i)
		fmt.Fprintf(fileTheor, "%5f %5f %5f\n", theorAvgUsers, theorAvgTime, i)
	}
}

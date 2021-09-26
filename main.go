package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blackNIKboard/queue-system-modeling/models"
	"github.com/blackNIKboard/queue-system-modeling/sync_system"
)

func main() {
	ss := sync_system.NewSyncSystem()

	//spew.Dump(poisson.GeneratePoissonProcess(0.05))

	//return

	if err := ss.Start(0); err != nil {
		return
	}

	for i := 0; i < 30; i++ {
		if err := ss.SendRequest(models.Request{
			IsFinished: false,
			AppendTime: time.Now(),
			EndTime:    time.Time{},
		}); err != nil {
			return
		}

		time.Sleep(time.Second / 2)
	}

	fmt.Println(ss.GetAvgTime())

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
}

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/blackNIKboard/queue-system-modeling/sync_system"
)

func main() {
	ss := sync_system.NewSyncSystem()

	if err := ss.Start(3); err != nil {
		return
	}

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
}

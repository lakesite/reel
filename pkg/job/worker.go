package job

import (
	"fmt"

	"github.com/lakesite/ls-governor"

	"github.com/lakesite/reel/pkg/reel"
)

// Queue for reel jobs.
var ReelQueue = make(chan ReelJob)

// ReelJob contains the information required to submit a job
// to a worker for processing.
type ReelJob struct {
	App string
	Source string
	Gapi *governor.API
}

type ReelWorker struct {}

func (rw *ReelWorker) Start() {
	go func() {
		for {
			work := <- ReelQueue
			fmt.Printf("Received a work request for app: %s using source: %s\n", work.App, work.Source)
			reel.Rewind(work.App, work.Source, work.Gapi)
		}
	}()
}

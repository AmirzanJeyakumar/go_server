package handlers

import (
	"fmt"
	"go_server/job_queue"
	"log"
	"net/http"
	"time"
)

func SubmitJobHandler(jq *job_queue.JobQueue, jobs int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for i := 1; i <= jobs; i++ {
			job := job_queue.Job{
				ID: fmt.Sprintf("job-%d", i),
				Execute: func(id int) func() {
					return func() {
						log.Printf("Executing job %d\n", id)
						time.Sleep(2 * time.Second)
						log.Printf("Job %d finished\n", id)
					}
				}(i),
			}
			err := jq.SubmitJob(job)
			if err != nil {
				fmt.Fprintf(w, "Failed to submit job: %s\n", err.Error())
				continue
			}
			fmt.Fprintf(w, "Job submitted: %s\n", job.ID)
		}
	}
}

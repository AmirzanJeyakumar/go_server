package main

import (
	"fmt"
	"go_server/job_queue"
	"log"
	"net/http"
	"time"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

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

func main() {
	jobQueue := job_queue.NewJobQueue(10)
	go jobQueue.Run()

	http.HandleFunc("/helloworld", helloWorld)
	http.HandleFunc("/submit_job", SubmitJobHandler(jobQueue, 20))

	log.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Panic(err)
	}
}

package job_queue

import (
	"context"
	"fmt"
	"time"
)

type Job struct {
	ID      string
	Execute func()
}

type JobQueue struct {
	jobs    chan Job
	workers int
}

func NewJobQueue(capacity int) *JobQueue {
	return &JobQueue{
		jobs:    make(chan Job, capacity),
		workers: capacity,
	}
}

func (jq *JobQueue) SubmitJob(job Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// If the job queue is full, we will wait until it is accepted, unless the context timed out
	select {
	case jq.jobs <- job:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("job queue is full, job %s was not submitted", job.ID)
	}
}

// Run starts goroutines for all workers
func (jq *JobQueue) Run() {
	for i := 0; i < jq.workers; i++ {
		go jq.worker()
	}
}

// worker will block and wait for a new job to come through the channel. This allows the workers to wait and continously process jobs from the channel as they become available
func (jq *JobQueue) worker() {
	for job := range jq.jobs {
		jq.processJob(job)
	}
}

func (jq *JobQueue) processJob(job Job) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan struct{})

	go func() {
		job.Execute()
		close(done)
	}()

	select {
	case <-ctx.Done():
		fmt.Printf("Job %s timed out\n", job.ID)
	case <-done:
		fmt.Printf("Job %s completed\n", job.ID)
	}
}

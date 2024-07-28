package main

import (
	"go_server/handlers"
	"go_server/job_queue"
	"log"
	"net/http"
	"os"
)

func main() {
	jobQueue := job_queue.NewJobQueue(10)
	go jobQueue.Run()

	memoryStore := handlers.NewMemoryStore()

	connStr := os.Getenv("POSTGRES_CONN_STR")
	postgresStore, err := handlers.NewPostgresStore(connStr)
	if err != nil {
		log.Fatalf("Could not create PostgreSQL store: %s\n", err)
	}

	http.HandleFunc("/helloworld", handlers.HelloWorld)
	http.HandleFunc("/submit_job", handlers.SubmitJobHandler(jobQueue, 20))
	http.HandleFunc("/memory/get", memoryStore.GetHandler)
	http.HandleFunc("/memory/put", memoryStore.PutHandler)
	http.HandleFunc("/postgres/get", postgresStore.GetHandler)
	http.HandleFunc("/postgres/put", postgresStore.PutHandler)

	log.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Panic(err)
	}
}

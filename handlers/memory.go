package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type MemoryStore struct {
	sync.RWMutex
	data map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]string),
	}
}

func (ms *MemoryStore) PutHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		log.Printf("PUT /memory/put - key is required - %v", time.Since(start))
		return
	}

	var value map[string]string
	if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		log.Printf("PUT /memory/put - invalid JSON - %v", time.Since(start))
		return
	}

	ms.Lock()
	ms.data[key] = value["value"]
	ms.Unlock()

	log.Printf("PUT /memory/put for key %s completed in %v", key, time.Since(start))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "data stored successfully\n")
}

func (ms *MemoryStore) GetHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		log.Printf("GET /memory/get - key is required - %v", time.Since(start))
		return
	}

	ms.RLock()
	value, ok := ms.data[key]
	ms.RUnlock()

	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		log.Printf("GET /memory/get - key not found - %v", time.Since(start))
		return
	}

	log.Printf("GET /memory/get for key %s completed in %v", key, time.Since(start))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"value": value})
}

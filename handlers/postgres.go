package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS my_table (
			key TEXT PRIMARY KEY,
			value JSONB
		)
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to create table: %w", err)
	}

	return &PostgresStore{DB: db}, nil
}

func (ps *PostgresStore) PutHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "error converting data to JSON", http.StatusInternalServerError)
		return
	}

	_, err = ps.DB.Exec("INSERT INTO my_table (key, value) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value", key, string(jsonData))
	if err != nil {
		http.Error(w, "error inserting data", http.StatusInternalServerError)
		return
	}

	log.Printf("PUT /postgres/put for key %s completed in %v", key, time.Since(start))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "data stored successfully\n")
}

func (ps *PostgresStore) GetHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		log.Printf("GET /postgres/get - key is required - %v", time.Since(start))
		return
	}

	var value string
	err := ps.DB.QueryRow("SELECT value FROM my_table WHERE key = $1", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "key not found", http.StatusNotFound)
			log.Printf("GET /postgres/get - key not found - %v", time.Since(start))
		} else {
			http.Error(w, "error retrieving data", http.StatusInternalServerError)
			log.Printf("GET /postgres/get - error retrieving data - %v", time.Since(start))
		}
		return
	}

	log.Printf("GET /postgres/get for key %s completed in %v", key, time.Since(start))

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s \n", value)
}

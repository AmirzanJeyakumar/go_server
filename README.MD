# notes

## concepts

- web server
- simple job queue to practice goroutines, channels and contexts
- RESTful APIs
- docker
- in-memory + persistent storage (postgres)

### example commands

build docker container:
`sudo docker-compose up --build`

```
1.  curl -X PUT -d '{"value": "Hello, Memory!"}' http://localhost:8080/memory/put?key=greeting

2.  curl http://localhost:8080/memory/get?key=greeting

3.  curl -X PUT -d '{"value": "Hello, Postgres!"}' http://localhost:8080/postgres/put?key=greeting

4.  curl http://localhost:8080/postgres/get?key=greeting
```

output:
```
1.  data stored successfully
    PUT /memory/put for key greeting completed in 33.18µs

2.  {"value":"Hello, Memory!"}
    GET /memory/get for key greeting completed in 1.58µs

3.  data stored successfully
    PUT /postgres/put for key greeting completed in 11.224253ms

4.  {"value": "Hello, Postgres!"}
    GET /postgres/get for key greeting completed in 625.46µs
```

**interpretation**
- Performance: In-memory operations are faster due to their direct nature and lack of I/O operations. Database operations, while slower, offer persistence and scalability.
- Latency Differences: The latency difference between in-memory and database operations highlights the trade-offs between speed and durability. In-memory stores are ideal for temporary or high-speed needs, while databases are suited for reliable, persistent storage


### things i found interesting

**Channels**
- Ranging over channels in Go is a convenient way to continuously receive values from a channel until it is closed. This is particularly useful in concurrent programming when you want to process a series of values sent 
  from multiple goroutines.



**sync.Mutex**
- A Mutex is a mutual exclusion lock.
- Only one goroutine can hold the lock at a time, whether it is for reading or writing.
- All other goroutines attempting to acquire the lock will be blocked until the lock is released.



**sync.RWMutex**
- A RWMutex is a reader/writer mutual exclusion lock.
- It has two types of locks: read locks (RLock) and write locks (Lock).
- Multiple goroutines can hold a read lock simultaneously, as long as there are no write locks held.
- Only one goroutine can hold the write lock, and no other locks (read or write) can be held at the same time.
- This allows multiple readers to access the resource concurrently, increasing the efficiency of read-heavy workloads.



**Benefits of sync.RWMutex**
- Improved Concurrency for Read-Heavy Workloads: When there are multiple readers and infrequent writers, RWMutex allows multiple readers to access the resource simultaneously, improving throughput and reducing contention.
- Separation of Read and Write Locks: By distinguishing between read and write operations, RWMutex provides a more fine-grained locking mechanism, which can lead to better performance in scenarios where reads vastly     
  outnumber writes.



**SQL stuff**

`ON CONFLICT (key)`

- This clause specifies how to handle a conflict that occurs when an attempt is made to insert a row with a key value that already exists in the table.
- In PostgreSQL, ON CONFLICT allows you to define an action to take if the insert would violate a unique constraint (in this case, a conflict on the key).

`DO UPDATE SET value = EXCLUDED.value`

- If a conflict occurs (i.e., a row with the same key already exists), this part of the query specifies that PostgreSQL should update the existing row.
- EXCLUDED is a special table reference used in the ON CONFLICT clause. It represents the row that was attempted to be inserted but caused the conflict.
- value = EXCLUDED.value means that the value column of the existing row should be updated to the value from the row that caused the conflict.



**request stuff**

```
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]string{"value": value})
```

- The Content-Type header informs the client about the type of content that is being returned. Setting this to application/json tells the client (such as a web browser or a REST client) that the response body contains 
  JSON data.
- The Encode method writes the serialized JSON directly to the response body. This means the client will receive the data as a JSON object, which is a standard format for RESTful APIs.

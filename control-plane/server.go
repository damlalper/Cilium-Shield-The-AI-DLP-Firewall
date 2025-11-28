
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Event represents a redaction event.
type Event struct {
	Timestamp      time.Time `json:"timestamp"`
	SourcePodIP    string    `json:"source_pod_ip"`
	DestinationURL string    `json:"destination_url"`
	RedactedType   string    `json:"redacted_type"`
}

// eventStore is an in-memory store for redaction events.
// It uses a sync.RWMutex to ensure thread-safe access to the events slice.
// This is crucial because we will have multiple goroutines writing to the store concurrently.
type eventStore struct {
	sync.RWMutex
	events []Event
}

// NewEventStore creates a new event store.
func NewEventStore() *eventStore {
	return &eventStore{
		events: make([]Event, 0),
	}
}

// Add adds a new event to the store.
// It uses a write lock (Lock()) to ensure that only one goroutine can modify the slice at a time.
func (es *eventStore) Add(event Event) {
	es.Lock()
	defer es.Unlock()
	es.events = append(es.events, event)
}

// List returns all events from the store.
// It uses a read lock (RLock()) to allow multiple goroutines to read the slice concurrently.
// This is more performant than using a write lock if we have many clients reading the events.
func (es *eventStore) List() []Event {
	es.RLock()
	defer es.RUnlock()
	// Return a copy to avoid race conditions where the caller modifies the slice.
	eventsCopy := make([]Event, len(es.events))
	copy(eventsCopy, es.events)
	return eventsCopy
}

var store = NewEventStore()

// eventHandler handles the /api/v1/events endpoint.
func eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// We process each event in a separate goroutine.
	// This makes the handler non-blocking and allows the server to handle a high volume of events
	// from the Wasm filters without slowing them down. The Wasm filter can send the event and
	// immediately continue processing the request.
	go func(e Event) {
		// In a real-world scenario, you might want to do more complex processing here,
		// such as enriching the event with data from other sources,
		// sending it to a message queue (e.g., Kafka, RabbitMQ), or storing it in a persistent database.
		log.Printf("Received event: Source: %s, Destination: %s, Type: %s", e.SourcePodIP, e.DestinationURL, e.RedactedType)
		store.Add(e)
	}(event)

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Event received")
}

// listEventsHandler handles the /api/v1/events/list endpoint.
func listEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.List())
}

func main() {
	http.HandleFunc("/api/v1/events", eventHandler)
	http.HandleFunc("/api/v1/events/list", listEventsHandler)

	log.Println("Cilium-Shield Observer started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}


package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestEventHandler(t *testing.T) {
	// Reset the store for each test
	store = NewEventStore()

	t.Run("valid post request", func(t *testing.T) {
		event := Event{
			Timestamp:      time.Now(),
			SourcePodIP:    "10.0.0.1",
			DestinationURL: "http://api.openai.com",
			RedactedType:   "CREDIT_CARD",
		}
		body, _ := json.Marshal(event)
		req, err := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusAccepted {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusAccepted)
		}

		// Allow time for the goroutine to process the event
		time.Sleep(10 * time.Millisecond)

		events := store.List()
		if len(events) != 1 {
			t.Errorf("expected 1 event in store, got %d", len(events))
		}
		if events[0].SourcePodIP != "10.0.0.1" {
			t.Errorf("expected source IP to be 10.0.0.1, got %s", events[0].SourcePodIP)
		}
	})

	t.Run("invalid request method", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/events", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusMethodNotAllowed)
		}
	})

	t.Run("invalid request body", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer([]byte("invalid json")))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
}

func TestListEventsHandler(t *testing.T) {
	t.Run("with events in store", func(t *testing.T) {
		// Reset the store
		store = NewEventStore()
		store.Add(Event{SourcePodIP: "10.0.0.1"})
		store.Add(Event{SourcePodIP: "10.0.0.2"})

		req, err := http.NewRequest("GET", "/api/v1/events/list", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(listEventsHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var events []Event
		err = json.Unmarshal(rr.Body.Bytes(), &events)
		if err != nil {
			t.Fatalf("could not parse response body: %s", err)
		}

		if len(events) != 2 {
			t.Errorf("expected 2 events, got %d", len(events))
		}
	})

	t.Run("with empty store", func(t *testing.T) {
		// Reset the store
		store = NewEventStore()

		req, err := http.NewRequest("GET", "/api/v1/events/list", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(listEventsHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var events []Event
		err = json.Unmarshal(rr.Body.Bytes(), &events)
		if err != nil {
			t.Fatalf("could not parse response body: %s", err)
		}

		if len(events) != 0 {
			t.Errorf("expected 0 events, got %d", len(events))
		}
	})
}

func TestEventStore_ListReturnsCopy(t *testing.T) {
	store := NewEventStore()
	store.Add(Event{SourcePodIP: "10.0.0.1"})

	events1 := store.List()
	if len(events1) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events1))
	}

	// Modify the returned slice
	events1[0].SourcePodIP = "modified"

	// Get the list again and check if it was modified
	events2 := store.List()
	if len(events2) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events2))
	}

	if events2[0].SourcePodIP == "modified" {
		t.Errorf("List method returned a reference to the internal slice, modifications should not be reflected")
	}
}

func TestEventStoreConcurrency(t *testing.T) {
	store := NewEventStore()
	var wg sync.WaitGroup
	numGoroutines := 100
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			store.Add(Event{SourcePodIP: "10.0.0." + string(i)})
		}(i)
	}

	wg.Wait()

	events := store.List()
	if len(events) != numGoroutines {
		t.Errorf("expected %d events in store, got %d", numGoroutines, len(events))
	}
}

package server

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// visitor holds data for each request (visitor)
type visitor struct {
	// limiter  *rate.Limiter
	lastSeen time.Time
	// time to send second request
	deadline time.Duration
	// Example for future use
	body string
}

// NewVisitor creates new visitor struct with timestamp
func newVisitor(maxDeadline int) *visitor {
	v := &visitor{lastSeen: time.Now()}
	v.deadline = (time.Duration(rand.Intn(maxDeadline) + 1)) * time.Second
	return v
}

// startTimer is called when first request is received
func (s *server) startTimer() {
	s.mu.Lock()
	s.counter++
	// Prevent filling memory with unclosed timers
	// also prevents integer overflow
	if s.counter >= s.config.ResetCounter {
		s.counter = 1
	}
	s.mu.Unlock()
	s.visitors[s.counter] = newVisitor(s.config.MaxDeadline)
	if s.config.Logging {
		log.Printf("ID: [%d], COUNT: [%d]\n", s.counter, len(s.visitors))
	}
}

// stopTimer is called when second request is received
func (s *server) stopTimer(query string) (time.Duration, time.Duration, error) {
	now := time.Now()
	id, err := strconv.Atoi(query)
	if err != nil {
		return 0, 0, errors.New("Parsing error: " + query)
	}
	if _, ok := s.visitors[id]; !ok {
		return 0, 0, errors.New("[Key] error")
	}
	delta := now.Sub(s.visitors[id].lastSeen)
	if s.config.Logging {
		log.Println("Time:", delta)
	}
	timeLimit := s.visitors[id].deadline
	delete(s.visitors, id)
	return delta, timeLimit, nil
}

// CleanVisitors cleans memory
func (s *server) CleanVisitors() {
	for {
		time.Sleep(s.config.WatchInterval)
		for id, v := range s.visitors {
			if time.Since(v.lastSeen) > s.config.TimeAlive {
				delete(s.visitors, id)
				if s.config.Logging {
					log.Println("Visitor deleted!")
				}
			}
		}
	}
}

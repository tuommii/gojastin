package server

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// visitor holds data for each request
type visitor struct {
	id int
	// limiter  *rate.Limiter
	lastSeen time.Time
	// time to send second request
	deadline time.Duration
	// Example for future use
	body string
}

// NewVisitor creates new visitor struct with timestamp
func newVisitor(deadline int) *visitor {
	v := &visitor{lastSeen: time.Now()}
	v.deadline = (time.Duration(rand.Intn(deadline) + 1)) * time.Second
	return v
}

// startTimer is called when first request is received
func (s *server) startTimer() {
	s.mu.Lock()
	s.counter++
	// Prevent filling memory with unclosed timers
	// also prevents integer overflow
	if s.counter >= s.config.MaxVisitors {
		s.counter = 0
	}
	s.mu.Unlock()
	s.visitors[s.counter] = newVisitor(s.config.Deadline)
	s.pool.Put(s.visitors)
	if s.config.Logging {
		log.Printf("id: [%d], count: [%d]\n", s.counter, len(s.visitors))
	}
}

// stopTimer is called when second request is received
func (s *server) stopTimer(query string) (time.Duration, time.Duration, error) {
	now := time.Now()
	id, err := strconv.Atoi(query)
	if err != nil {
		return 0, 0, errors.New("parsing error: " + query)
	}
	if _, ok := s.visitors[id]; !ok {
		return 0, 0, errors.New("visitor doesn't exist")
	}

	// m := s.pool.Get().(map[int]*visitor)
	delta := now.Sub(s.visitors[id].lastSeen)
	if s.config.Logging {
		log.Println("time:", delta)
	}
	timeLimit := s.visitors[id].deadline
	delete(s.visitors, id)
	s.pool.Put(s.visitors)
	return delta, timeLimit, nil
}

// CleanVisitors cleans memory
func (s *server) CleanVisitors() {
	for {
		time.Sleep(s.config.RemoveInterval)
		for id, v := range s.visitors {
			if time.Since(v.lastSeen) > s.config.VisitorAlive {
				delete(s.visitors, id)
				if s.config.Logging {
					log.Println("visitor deleted!")
				}
			}
		}
		s.pool.Put(s.visitors)
	}
}

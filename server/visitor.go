package server

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const timeAlive = time.Second * 30

// visitor holds data for each request (visitor)
type visitor struct {
	// limiter  *rate.Limiter
	lastSeen time.Time
	deadline time.Duration
	// Example for future use
	body string
}

// NewVisitor creates new visitor struct with timestamp
func newVisitor() *visitor {
	v := &visitor{lastSeen: time.Now()}
	v.deadline = (time.Duration(rand.Intn(10) + 1)) * time.Second
	return v
}

// startTimer is called when first request is received
func (s *server) startTimer() {
	s.mu.Lock()
	s.counter++
	// Prevent filling memory with unclosed timers
	// also prevents integer overflow
	if s.counter >= 1000 {
		s.counter = 1
	}
	s.mu.Unlock()
	s.visitors[s.counter] = newVisitor()
	if s.logging {
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
	if s.logging {
		log.Println("Time:", delta)
	}
	timeLimit := s.visitors[id].deadline
	delete(s.visitors, id)
	return delta, timeLimit, nil
}

// CleanVisitors cleans memory
func (s *server) CleanVisitors() {
	for {
		time.Sleep(time.Minute * 2)
		for ip, v := range s.visitors {
			if time.Since(v.lastSeen) > timeAlive {
				delete(s.visitors, ip)
			}
		}
	}
}

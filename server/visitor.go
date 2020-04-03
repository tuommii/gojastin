package server

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const cleanUpTime = time.Second * 30

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

// startTiming is called when first request is received
func (s *server) startTiming() {
	s.mu.Lock()
	s.counter++
	// Prevent filling memory with unclosed timers
	// also prevents integer going over max value
	if s.counter >= 1000 {
		s.counter = 1
	}
	s.mu.Unlock()
	s.visitors[s.counter] = newVisitor()
	log.Printf("ID: [%d], COUNT: [%d]\n", s.counter, len(s.visitors))
}

// stopTiming is called when second request is received
func (s *server) stopTiming(query string) (time.Duration, time.Duration, error) {
	now := time.Now()
	id, err := strconv.Atoi(query)
	if err != nil {
		return 0, 0, errors.New("Parsing error: " + query)
	}
	if _, ok := s.visitors[id]; !ok {
		return 0, 0, errors.New("[Key] error")
	}
	delta := now.Sub(s.visitors[id].lastSeen)
	log.Println("Time:", delta)
	timeLimit := s.visitors[id].deadline
	delete(s.visitors, id)
	return delta, timeLimit, nil
}

// CleanVisitors checks every minute inactive visitors and delete's them
func (s *server) CleanVisitors() {
	for {
		time.Sleep(time.Minute * 2)
		for ip, v := range s.visitors {
			if time.Since(v.lastSeen) > cleanUpTime {
				delete(s.visitors, ip)
			}
		}
	}
}

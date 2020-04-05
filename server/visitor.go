package server

import (
	"errors"
	"fmt"
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
	// s.visitors[s.counter] = newVisitor(s.config.Deadline)
	v := s.pool.Get().(map[int]*visitor)
	if s.config.Logging {
		fmt.Printf("before: %+v\n", v)
	}
	v[s.counter] = newVisitor(s.config.Deadline)
	s.pool.Put(v)
	if s.config.Logging {
		fmt.Printf("after: %+v\n", v)
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

	v := s.pool.Get().(map[int]*visitor)
	if _, ok := v[id]; !ok {
		return 0, 0, errors.New("visitor doesn't exist")
	}

	delta := now.Sub(v[id].lastSeen)
	if s.config.Logging {
		log.Println("time:", delta)
	}
	timeLimit := v[id].deadline
	delete(v, id)
	s.pool.Put(v)
	return delta, timeLimit, nil
}

// CleanVisitors cleans memory
func (s *server) CleanVisitors() {
	for {
		time.Sleep(s.config.RemoveInterval)
		vis := s.pool.Get().(map[int]*visitor)
		for id, v := range vis {
			if time.Since(v.lastSeen) > s.config.VisitorAlive {
				delete(vis, id)
				if s.config.Logging {
					log.Println("visitor deleted!")
				}
			}
		}
		s.pool.Put(vis)
	}
}

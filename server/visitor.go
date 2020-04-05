package server

import (
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
	body int
}

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
	v := s.pool.Get().(*visitor)
	v.body = s.counter
	s.pool.Put(v)
	s.visitors[s.counter] = v
}

// stopTimer is called when second request is received
func (s *server) stopTimer(query string) (time.Duration, *visitor) {
	now := time.Now()
	id, err := strconv.Atoi(query)
	if err != nil {
		return 0, nil
	}
	if _, ok := s.visitors[id]; !ok {
		return 0, nil
	}
	delta := now.Sub(s.visitors[id].lastSeen)
	if s.config.Logging {
		fmt.Printf("%+v, %s\n", s.visitors[id], delta)
	}
	return delta, s.visitors[id]
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
	}
}

package server

import (
	"log"
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

	// Here time gets added to visitor
	v := s.pool.Get().(*visitor)
	v.id = s.counter
	s.pool.Put(v)
	s.visitors[s.counter] = v
}

// stopTimer is called when second request is received
func (s *server) stopTimer(query string) (time.Duration, *visitor) {
	now := time.Now()
	id, err := parseQuery(query)

	if err != nil {
		log.Println(err)
		return 0, nil
	}
	if _, exist := s.visitors[id]; !exist {
		return 0, nil
	}

	delta := now.Sub(s.visitors[id].lastSeen)
	return delta, s.visitors[id]
}

// Parse id from query/url
func parseQuery(query string) (int, error) {
	id, err := strconv.Atoi(query)
	if err != nil {
		return 0, err
	}
	return id, nil
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

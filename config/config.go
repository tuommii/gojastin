package config

import "time"

// Config for server
type Config struct {
	// (sec) older than this gets removed in background routine
	VisitorAlive time.Duration
	// Every n look old visitors to be deleted (sec)
	RemoveInterval time.Duration
	// Maximux number (int) to send second request (sec)
	Deadline int
	// after this counter is set back to 1
	MaxVisitors int
	// Is logging enabled, default true
	Logging bool
}

// New returns new config
func New() *Config {
	c := &Config{
		VisitorAlive:   time.Second * 30,
		RemoveInterval: time.Second * 120,
		Deadline:       7,
		MaxVisitors:    1000,
		Logging:        true,
	}
	if c.RemoveInterval < c.VisitorAlive {
		c.RemoveInterval = c.VisitorAlive
	}
	return c
}

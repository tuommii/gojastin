package config

import "time"

// Every n:th (sec) look old visitors to be deleted
const removeInterval = time.Second * 120

// older than this gets removed in background routine
const alive = time.Second * 30

// after this counter is set back to 1
const reset = 1000

// maximux time (sec) to send second request
const deadline = 7

// Config for server
type Config struct {
	// older than this gets removed in background routine
	Alive time.Duration
	// Every n look old visitors to be deleted (sec)
	RemoveInterval time.Duration
	// maximux time to send second request (sec)
	Deadline int
	// after this counter is set back to 1
	Reset int
	// Is logging enabled, default true
	Logging bool
}

// New returns new config
func New() *Config {
	c := &Config{
		Alive:          alive,
		Reset:          reset,
		Deadline:       deadline,
		RemoveInterval: removeInterval,
		Logging:        true,
	}
	if c.RemoveInterval < alive {
		c.RemoveInterval = alive
	}
	return c
}

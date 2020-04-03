package config

import "time"

// Every n:th (sec) look old visitors to be deleted
const watchInterval = time.Second * 120

// older than this gets removed in background routine
const timeAlive = time.Second * 30

// after this counter is set back to 1
const resetCounter = 1000

// maximux time (sec) to send second request
const maxDeadline = 7

// Config for server
type Config struct {
	// older than this gets removed in background routine
	TimeAlive time.Duration
	// Every n look old visitors to be deleted (sec)
	WatchInterval time.Duration
	// maximux time to send second request (sec)
	MaxDeadline int
	// after this counter is set back to 1
	ResetCounter int
	// Is logging enabled, default true
	Logging bool
}

// New returns new config
func New() *Config {
	c := &Config{
		TimeAlive:     timeAlive,
		ResetCounter:  resetCounter,
		MaxDeadline:   maxDeadline,
		WatchInterval: watchInterval,
		Logging:       true,
	}
	if c.WatchInterval < timeAlive {
		c.WatchInterval = timeAlive
	}
	return c
}

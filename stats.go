package main

import (
	"fmt"
	"sync"
	"time"
)

type Stats struct {
	nBytes     uint64
	timeStart  time.Time
	timeStop   time.Time
	timePause  time.Time
	timePaused time.Duration

	lock *sync.RWMutex
}

// NewStats creates a new Stats
func NewStats() *Stats {
	return &Stats{
		lock: &sync.RWMutex{},
	}
}

// Start stores the "start" timestamp
func (s *Stats) Start() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.timeStart.IsZero() {
		s.timeStart = time.Now()
	} else if !s.timePause.IsZero() {
		s.timePaused += time.Since(s.timePause)
		// Reset
		s.timePause = time.Time{}
	}
}

// Pause stores an interruption timestamp
func (s *Stats) Pause() {
	s.lock.RLock()

	if s.timeStart.IsZero() || !s.timeStop.IsZero() {
		// Can't stop if not started, or if stopped
		s.lock.RUnlock()
		return
	}
	s.lock.RUnlock()

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.timePause.IsZero() {
		s.timePause = time.Now()
	}
}

// Stop stores the stop timestamp
func (s *Stats) Stop() {
	s.lock.RLock()

	if s.timeStart.IsZero() {
		// Can't stop if not started
		s.lock.RUnlock()
		return
	}
	s.lock.RUnlock()

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.timeStop.IsZero() {
		s.timeStop = time.Now()
	}
}

// Speed get the IO speed in MB/s
func (s *Stats) Speed() float64 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return float64(s.nBytes) / 1024 / 1024 / s.Duration().Seconds()
}

// Duration get the stop start duration
func (s *Stats) Duration() time.Duration {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.timeStart.IsZero() {
		return 0
	} else if s.timeStop.IsZero() {
		return time.Since(s.timeStart) - s.timePaused
	}
	return s.timeStop.Sub(s.timeStart) - s.timePaused
}

// Bytes get the stored number of bytes
func (s *Stats) Bytes() uint64 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.nBytes
}

// AddBytes increase the nbBytes counter
func (s *Stats) AddBytes(c uint64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.nBytes += c
}

func (s *Stats) String() string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return fmt.Sprintf("%v bytes | %-v | %0.4f MB/s", s.Bytes(), s.Duration(), s.Speed())
}

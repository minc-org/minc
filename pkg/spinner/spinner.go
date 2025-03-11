package spinner

import (
	"fmt"
	"sync"
	"time"
)

type Spinner struct {
	mu       *sync.RWMutex
	Delay    time.Duration
	active   bool
	stopChan chan bool
}

func New(d time.Duration) *Spinner {
	return &Spinner{
		Delay:    d,
		mu:       &sync.RWMutex{},
		active:   false,
		stopChan: make(chan bool, 1),
	}
}

func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	go func() {
		for {
			select {
			case <-s.stopChan:
				return
			default:
				s.mu.RLock()
				if !s.active {
					s.mu.RUnlock()
					return
				}
				s.mu.RUnlock()

				for _, r := range `|/-\` {
					fmt.Printf("\r%c", r)
					time.Sleep(s.Delay)

					s.mu.RLock()
					if !s.active {
						s.mu.RUnlock()
						return
					}
					s.mu.RUnlock()
				}
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.active = false
		select {
		case s.stopChan <- true:
		default:
		}
	}
}

func (s *Spinner) Active() bool {
	return s.active
}

package auth

import (
	"sync"
	"time"
)

type LoginLimiter struct {
	mu        sync.Mutex
	max       int
	window    time.Duration
	now       func() time.Time
	failures  map[string][]time.Time
	lastSweep time.Time
}

func NewLoginLimiter(max int, window time.Duration) *LoginLimiter {
	return &LoginLimiter{max: max, window: window, now: time.Now, failures: map[string][]time.Time{}}
}

func (l *LoginLimiter) recent(key string) []time.Time {
	cutoff := l.now().Add(-l.window)
	kept := l.failures[key][:0]
	for _, t := range l.failures[key] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) == 0 {
		delete(l.failures, key)
	} else {
		l.failures[key] = kept
	}
	return kept
}

func (l *LoginLimiter) Blocked(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.recent(key)) >= l.max
}

func (l *LoginLimiter) Fail(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sweep()
	l.failures[key] = append(l.recent(key), l.now())
}

func (l *LoginLimiter) sweep() {
	if l.now().Sub(l.lastSweep) < l.window {
		return
	}
	l.lastSweep = l.now()
	for key := range l.failures {
		l.recent(key)
	}
}

func (l *LoginLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.failures, key)
}

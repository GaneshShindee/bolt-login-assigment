package auth

import (
	"sync"
	"testing"
	"time"
)

func newTestLimiter(max int, window time.Duration) (*LoginLimiter, *time.Time) {
	now := time.Unix(1_700_000_000, 0)
	l := NewLoginLimiter(max, window)
	l.now = func() time.Time { return now }
	return l, &now
}

func TestLimiter_BlocksAfterMaxFailures(t *testing.T) {
	l, _ := newTestLimiter(3, time.Minute)
	for i := range 3 {
		if l.Blocked("a@x.com") {
			t.Fatalf("blocked after only %d failures", i)
		}
		l.Fail("a@x.com")
	}
	if !l.Blocked("a@x.com") {
		t.Fatal("should be blocked after 3 failures")
	}
	if l.Blocked("b@x.com") {
		t.Fatal("other keys must not be affected")
	}
}

func TestLimiter_WindowSlides(t *testing.T) {
	l, now := newTestLimiter(2, time.Minute)
	l.Fail("a@x.com")
	l.Fail("a@x.com")
	if !l.Blocked("a@x.com") {
		t.Fatal("should be blocked")
	}
	*now = now.Add(61 * time.Second)
	if l.Blocked("a@x.com") {
		t.Fatal("old failures should expire after the window")
	}
}

func TestLimiter_ResetClears(t *testing.T) {
	l, _ := newTestLimiter(1, time.Minute)
	l.Fail("a@x.com")
	l.Reset("a@x.com")
	if l.Blocked("a@x.com") {
		t.Fatal("Reset should clear failures")
	}
}

func TestLimiter_ForgetsExpiredKeys(t *testing.T) {
	l, now := newTestLimiter(5, time.Minute)
	for _, key := range []string{"a@x.com", "b@x.com", "c@x.com"} {
		l.Fail(key) // each fails once and never retries
	}
	*now = now.Add(2 * time.Minute)
	l.Fail("d@x.com") // triggers a sweep
	if len(l.failures) != 1 {
		t.Fatalf("expired keys should be dropped, %d keys left: %v", len(l.failures), l.failures)
	}
}

func TestLimiter_ConcurrentSafe(t *testing.T) {
	l := NewLoginLimiter(1000, time.Minute)
	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() {
			l.Fail("a@x.com")
			l.Blocked("a@x.com")
		})
	}
	wg.Wait()
	if got := len(l.failures["a@x.com"]); got != 50 {
		t.Fatalf("got %d failures, want 50", got)
	}
}

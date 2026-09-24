package auth

import (
	"strings"
	"testing"
	"time"
)

func TestToken_RoundTrip(t *testing.T) {
	s := NewTokenSigner("secret", time.Hour)
	id, err := s.Verify(s.Issue(42))
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if id != 42 {
		t.Fatalf("got user id %d, want 42", id)
	}
}

func TestToken_Rejects(t *testing.T) {
	s := NewTokenSigner("secret", time.Hour)
	valid := s.Issue(42)
	parts := strings.Split(valid, ".")

	tests := map[string]string{
		"empty":            "",
		"garbage":          "not-a-token",
		"too many parts":   valid + ".extra",
		"tampered user id": "1." + parts[1] + "." + parts[2],
		"tampered expiry":  parts[0] + ".9999999999." + parts[2],
		"tampered sig":     parts[0] + "." + parts[1] + ".AAAA",
		"other secret":     NewTokenSigner("different", time.Hour).Issue(42),
	}
	for name, tok := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := s.Verify(tok); err == nil {
				t.Fatalf("token %q should be rejected", tok)
			}
		})
	}
}

func TestToken_Expires(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	s := NewTokenSigner("secret", time.Hour)
	s.now = func() time.Time { return now }
	tok := s.Issue(7)

	now = now.Add(59 * time.Minute)
	if _, err := s.Verify(tok); err != nil {
		t.Fatalf("token should still be valid before expiry: %v", err)
	}
	now = now.Add(2 * time.Minute)
	if _, err := s.Verify(tok); err == nil {
		t.Fatal("token should be rejected after expiry")
	}
}

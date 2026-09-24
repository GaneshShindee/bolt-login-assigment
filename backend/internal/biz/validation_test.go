package biz

import (
	"strings"
	"testing"
)

func TestValidEmail(t *testing.T) {
	valid := []string{
		"a@b.co",
		"test.user@example.com",
		"first+tag@sub.domain.org",
		"x_y-z%1@my-site.in",
	}
	invalid := []string{
		"",
		"plain",
		"a@",
		"@b.com",
		"a@b",   // no TLD: still typing
		"a@b.c", // 1-letter TLD
		"a@b.",
		"a b@c.com",
		"a@-b.com",
		"a@b-.com",
		"a@@b.com",
		strings.Repeat("a", 250) + "@b.com", // longer than 254
	}
	for _, e := range valid {
		if !validEmail(e) {
			t.Errorf("validEmail(%q) = false, want true", e)
		}
	}
	for _, e := range invalid {
		if validEmail(e) {
			t.Errorf("validEmail(%q) = true, want false", e)
		}
	}
}

func TestNormalizeEmail(t *testing.T) {
	if got := normalizeEmail("  Test.User@Example.COM \n"); got != "test.user@example.com" {
		t.Fatalf("got %q", got)
	}
}

func TestValidCode(t *testing.T) {
	for _, c := range []string{"000000", "123456", "999999"} {
		if !validCode(c) {
			t.Errorf("validCode(%q) = false, want true", c)
		}
	}
	for _, c := range []string{"", "12345", "1234567", "12a456", " 123456", "１２３４５６"} {
		if validCode(c) {
			t.Errorf("validCode(%q) = true, want false", c)
		}
	}
}

func TestValidPhone(t *testing.T) {
	for _, p := range []string{"9876543210", "6000000000", "7012345678", "8999999999"} {
		if !validPhone(p) {
			t.Errorf("validPhone(%q) = false, want true", p)
		}
	}
	for _, p := range []string{
		"", "987654321", "98765432101", // 9 and 11 digits
		"5876543210", "0987654321", // must start with 6–9
		"+919876543210", "98765 43210", "98765-43210", "98765abcde", // digits only
	} {
		if validPhone(p) {
			t.Errorf("validPhone(%q) = true, want false", p)
		}
	}
}

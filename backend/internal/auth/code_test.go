package auth

import (
	"regexp"
	"testing"
)

func TestGenerateCode_IsSixDigits(t *testing.T) {
	re := regexp.MustCompile(`^\d{6}$`)
	for range 1000 {
		code, err := GenerateCode()
		if err != nil {
			t.Fatalf("GenerateCode: %v", err)
		}
		if !re.MatchString(code) {
			t.Fatalf("code %q is not exactly 6 digits", code)
		}
	}
}

func TestGenerateCode_IsRandom(t *testing.T) {
	seen := map[string]bool{}
	for range 200 {
		code, _ := GenerateCode()
		seen[code] = true
	}
	// 200 draws from 1,000,000: collisions are possible but many repeats means it's not random.
	if len(seen) < 190 {
		t.Fatalf("expected ~200 distinct codes, got %d", len(seen))
	}
}

func TestHashAndCheckCode(t *testing.T) {
	hash, err := HashCode("042917")
	if err != nil {
		t.Fatalf("HashCode: %v", err)
	}
	if hash == "042917" {
		t.Fatal("hash must not equal the plain code")
	}
	if !CheckCode(hash, "042917") {
		t.Error("correct code should match")
	}
	if CheckCode(hash, "042918") {
		t.Error("wrong code should not match")
	}
	if CheckCode(hash, "42917") {
		t.Error("code without its leading zero should not match")
	}
}

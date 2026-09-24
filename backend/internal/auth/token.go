package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrBadToken = errors.New("invalid token")

type TokenSigner struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

func NewTokenSigner(secret string, ttl time.Duration) *TokenSigner {
	return &TokenSigner{secret: []byte(secret), ttl: ttl, now: time.Now}
}

func (t *TokenSigner) sign(payload string) string {
	mac := hmac.New(sha256.New, t.secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (t *TokenSigner) Issue(userID int) string {
	payload := fmt.Sprintf("%d.%d", userID, t.now().Add(t.ttl).Unix())
	return payload + "." + t.sign(payload)
}

func (t *TokenSigner) Verify(token string) (int, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, ErrBadToken
	}
	payload := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(t.sign(payload)), []byte(parts[2])) {
		return 0, ErrBadToken
	}
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || t.now().Unix() > exp {
		return 0, ErrBadToken
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, ErrBadToken
	}
	return id, nil
}

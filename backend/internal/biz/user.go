package biz

import (
	"context"
	"errors"
	"strings"

	"github.com/ganeshshinde/boltapp/backend/internal/auth"
	"github.com/ganeshshinde/boltapp/backend/internal/entity"
)

type UserBiz struct {
	repo    UserRepo
	tokens  *auth.TokenSigner
	limiter *auth.LoginLimiter
}

func NewUserBiz(repo UserRepo, tokens *auth.TokenSigner, limiter *auth.LoginLimiter) *UserBiz {
	return &UserBiz{repo: repo, tokens: tokens, limiter: limiter}
}

func (b *UserBiz) Register(ctx context.Context, email, firstName, lastName string) (entity.User, string, error) {
	u := entity.User{
		Email:     normalizeEmail(email),
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
	}
	switch {
	case !validEmail(u.Email):
		return entity.User{}, "", invalid("Please enter a valid email address.")
	case !validName(u.FirstName):
		return entity.User{}, "", invalid("Please enter your first name.")
	case !validName(u.LastName):
		return entity.User{}, "", invalid("Please enter your last name.")
	}

	code, err := auth.GenerateCode()
	if err != nil {
		return entity.User{}, "", err
	}
	if u.LoginCodeHash, err = auth.HashCode(code); err != nil {
		return entity.User{}, "", err
	}
	u, err = b.repo.Create(ctx, u)
	if errors.Is(err, ErrDuplicate) {
		return entity.User{}, "", ErrEmailTaken
	}
	if err != nil {
		return entity.User{}, "", err
	}
	return u, code, nil
}

func (b *UserBiz) Recognize(ctx context.Context, email string) (bool, error) {
	email = normalizeEmail(email)
	if !validEmail(email) {
		return false, invalid("Please enter a valid email address.")
	}
	_, err := b.repo.ByEmail(ctx, email)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (b *UserBiz) Login(ctx context.Context, email, code string) (entity.User, string, error) {
	email = normalizeEmail(email)
	code = strings.TrimSpace(code)
	if !validEmail(email) || !validCode(code) {
		return entity.User{}, "", invalid("Please enter your 6-digit code.")
	}
	if b.limiter.Blocked(email) {
		return entity.User{}, "", ErrTooManyAttempts
	}

	u, err := b.repo.ByEmail(ctx, email)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return entity.User{}, "", err
	}
	// Unknown email and wrong code fail the same way, and both count toward the limit.
	if errors.Is(err, ErrNotFound) || !auth.CheckCode(u.LoginCodeHash, code) {
		b.limiter.Fail(email)
		return entity.User{}, "", ErrInvalidCode
	}
	b.limiter.Reset(email)
	return u, b.tokens.Issue(u.ID), nil
}

func (b *UserBiz) Authenticate(ctx context.Context, token string) (entity.User, error) {
	id, err := b.tokens.Verify(token)
	if err != nil {
		return entity.User{}, ErrUnauthenticated
	}
	u, err := b.repo.ByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return entity.User{}, ErrUnauthenticated
	}
	return u, err
}

// Package repotest provides in-memory biz repositories, so tests run without Postgres.
package repotest

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/ganeshshinde/boltapp/backend/internal/biz"
	"github.com/ganeshshinde/boltapp/backend/internal/entity"
)

// Users is an in-memory biz.UserRepo. All is exported so tests can inspect it.
type Users struct {
	mu  sync.Mutex
	All []entity.User
}

func (r *Users) Create(_ context.Context, u entity.User) (entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.All {
		if existing.Email == u.Email {
			return entity.User{}, biz.ErrDuplicate
		}
	}
	u.ID = len(r.All) + 1
	r.All = append(r.All, u)
	return u, nil
}

func (r *Users) ByEmail(_ context.Context, email string) (entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.All {
		if u.Email == email {
			return u, nil
		}
	}
	return entity.User{}, biz.ErrNotFound
}

func (r *Users) ByID(_ context.Context, id int) (entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.All {
		if u.ID == id {
			return u, nil
		}
	}
	return entity.User{}, biz.ErrNotFound
}

type Checkouts struct {
	mu  sync.Mutex
	All []entity.Checkout
}

func (r *Checkouts) Create(_ context.Context, c entity.Checkout) (entity.Checkout, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c.ID = len(r.All) + 1
	c.CreatedAt = time.Now()
	r.All = append(r.All, c)
	return c, nil
}

func (r *Checkouts) RecentDetails(_ context.Context, userID, limit int) ([]entity.SavedDetails, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entity.SavedDetails
	// Newest first; keep the first (most recent) occurrence of each phone + address.
	for _, c := range slices.Backward(r.All) {
		if c.UserID == nil || *c.UserID != userID {
			continue
		}
		dup := slices.ContainsFunc(out, func(d entity.SavedDetails) bool {
			return d.Phone == c.Phone && d.Address == c.Address
		})
		if !dup {
			out = append(out, entity.SavedDetails{Phone: c.Phone, Address: c.Address, LastUsedAt: c.CreatedAt})
		}
		if len(out) == limit {
			break
		}
	}
	return out, nil
}

var (
	_ biz.UserRepo     = (*Users)(nil)
	_ biz.CheckoutRepo = (*Checkouts)(nil)
)

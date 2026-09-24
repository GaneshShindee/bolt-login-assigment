package biz

import (
	"context"

	"github.com/ganeshshinde/boltapp/backend/internal/entity"
)

// Repository interfaces live where they are used; the repository package (Postgres) and
// repotest (in-memory) implement them.

type UserRepo interface {
	// Create returns ErrDuplicate if the email exists.
	Create(ctx context.Context, u entity.User) (entity.User, error)
	// ByEmail returns ErrNotFound if no user has this email.
	ByEmail(ctx context.Context, email string) (entity.User, error)
	// ByID returns ErrNotFound if no user has this ID.
	ByID(ctx context.Context, id int) (entity.User, error)
}

type CheckoutRepo interface {
	Create(ctx context.Context, c entity.Checkout) (entity.Checkout, error)
	// RecentDetails returns up to limit distinct phone + address pairs, most recently used first.
	RecentDetails(ctx context.Context, userID, limit int) ([]entity.SavedDetails, error)
}

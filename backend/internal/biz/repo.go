package biz

import (
	"context"

	"github.com/ganeshshinde/boltapp/backend/internal/entity"
)

type UserRepo interface {
	Create(ctx context.Context, u entity.User) (entity.User, error)
	ByEmail(ctx context.Context, email string) (entity.User, error)
	ByID(ctx context.Context, id int) (entity.User, error)
}

type CheckoutRepo interface {
	Create(ctx context.Context, c entity.Checkout) (entity.Checkout, error)
	RecentDetails(ctx context.Context, userID, limit int) ([]entity.SavedDetails, error)
}

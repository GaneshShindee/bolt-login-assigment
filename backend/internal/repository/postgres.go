package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ganeshshinde/boltapp/backend/internal/biz"
	"github.com/ganeshshinde/boltapp/backend/internal/entity"
)

type DB struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, url string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	// Supabase's transaction pooler (port 6543) doesn't support prepared statements.
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	cfg.MaxConns = 5 // free-tier Postgres has a small connection limit

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &DB{pool: pool}, nil
}

func (db *DB) Close() { db.pool.Close() }

func (db *DB) Ping(ctx context.Context) error { return db.pool.Ping(ctx) }

type UserRepo struct{ db *DB }

func NewUserRepo(db *DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, u entity.User) (entity.User, error) {
	err := r.db.pool.QueryRow(ctx,
		`INSERT INTO users (email, first_name, last_name, login_code_hash)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		u.Email, u.FirstName, u.LastName, u.LoginCodeHash,
	).Scan(&u.ID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
		return entity.User{}, biz.ErrDuplicate
	}
	return u, err
}

func (r *UserRepo) ByEmail(ctx context.Context, email string) (entity.User, error) {
	return scanUser(r.db.pool.QueryRow(ctx,
		`SELECT id, email, first_name, last_name, login_code_hash FROM users WHERE email = $1`, email))
}

func (r *UserRepo) ByID(ctx context.Context, id int) (entity.User, error) {
	return scanUser(r.db.pool.QueryRow(ctx,
		`SELECT id, email, first_name, last_name, login_code_hash FROM users WHERE id = $1`, id))
}

func scanUser(row pgx.Row) (entity.User, error) {
	var u entity.User
	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.LoginCodeHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.User{}, biz.ErrNotFound
	}
	return u, err
}

type CheckoutRepo struct{ db *DB }

func NewCheckoutRepo(db *DB) *CheckoutRepo { return &CheckoutRepo{db: db} }

func (r *CheckoutRepo) Create(ctx context.Context, c entity.Checkout) (entity.Checkout, error) {
	a := c.Address
	err := r.db.pool.QueryRow(ctx,
		`INSERT INTO checkouts (user_id, email, phone, address_label, address_line1, address_line2, city, state, pincode)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at`,
		c.UserID, c.Email, c.Phone, a.Label, a.Line1, a.Line2, a.City, a.State, a.Pincode,
	).Scan(&c.ID, &c.CreatedAt)
	return c, err
}

func (r *CheckoutRepo) RecentDetails(ctx context.Context, userID, limit int) ([]entity.SavedDetails, error) {
	rows, err := r.db.pool.Query(ctx,
		`SELECT phone, address_label, address_line1, address_line2, city, state, pincode, MAX(created_at) AS last_used
		 FROM checkouts
		 WHERE user_id = $1
		 GROUP BY phone, address_label, address_line1, address_line2, city, state, pincode
		 ORDER BY last_used DESC
		 LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.SavedDetails, error) {
		var d entity.SavedDetails
		a := &d.Address
		err := row.Scan(&d.Phone, &a.Label, &a.Line1, &a.Line2, &a.City, &a.State, &a.Pincode, &d.LastUsedAt)
		return d, err
	})
}

var (
	_ biz.UserRepo     = (*UserRepo)(nil)
	_ biz.CheckoutRepo = (*CheckoutRepo)(nil)
)

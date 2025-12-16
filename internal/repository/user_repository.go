package repository

import (
	"context"
	"time"

	db "go-users-api/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository interface {
	Create(ctx context.Context, name string, dob time.Time) (db.User, error)
	GetByID(ctx context.Context, id int64) (db.User, error)
	List(ctx context.Context) ([]db.User, error)
	Update(ctx context.Context, id int64, name string, dob time.Time) (db.User, error)
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) UserRepository {
	return &userRepository{q: q}
}

func toPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func toInt32(id int64) int32 {
	return int32(id)
}

func (r *userRepository) Create(ctx context.Context, name string, dob time.Time) (db.User, error) {
	return r.q.CreateUser(ctx, db.CreateUserParams{
		Name: name,
		Dob:  toPgDate(dob),
	})
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (db.User, error) {
	return r.q.GetUser(ctx, toInt32(id))
}

func (r *userRepository) List(ctx context.Context) ([]db.User, error) {
	return r.q.ListUsers(ctx)
}

func (r *userRepository) Update(ctx context.Context, id int64, name string, dob time.Time) (db.User, error) {
	return r.q.UpdateUser(ctx, db.UpdateUserParams{
		Name: name,
		Dob:  toPgDate(dob),
		ID:   toInt32(id),
	})
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	return r.q.DeleteUser(ctx, toInt32(id))
}

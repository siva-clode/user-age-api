package repository

import (
	"context"
	"database/sql"
	"time"

	db "github.com/bugude99/user-age-api/db/sqlc"
)

type UserRepo struct {
	q *db.Queries
}

func NewUserRepo(dbConn *sql.DB) *UserRepo {
	return &UserRepo{
		q: db.New(dbConn),
	}
}

func (r *UserRepo) Create(ctx context.Context, name string, dob time.Time) (db.User, error) {
	return r.q.CreateUser(ctx, db.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (db.User, error) {
	return r.q.GetUserByID(ctx, int32(id))
}

func (r *UserRepo) Update(ctx context.Context, id int, name string, dob time.Time) (db.User, error) {
	return r.q.UpdateUser(ctx, db.UpdateUserParams{
		ID:   int32(id),
		Name: name,
		Dob:  dob,
	})
}

func (r *UserRepo) Delete(ctx context.Context, id int) error {
	return r.q.DeleteUser(ctx, int32(id))
}

func (r *UserRepo) List(ctx context.Context, limit int, offset int) ([]db.User, error) {
	users, err := r.q.ListUsers(ctx, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	return users, err
}

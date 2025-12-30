package repository

import (
	"context"
	"go-shopping-cart/internal/db/sqlc"
)

type SqlUserRepository struct {
	db sqlc.Querier
}

func NewSqlUserRepository(db sqlc.Querier) UserRepository {
	return &SqlUserRepository{
		db: db,
	}
}

func (ur *SqlUserRepository) FindAll() {}

func (ur *SqlUserRepository) Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.db.CreateUser(ctx, userParams)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SqlUserRepository) FindByUUID(uuid string) {}

func (ur *SqlUserRepository) Update(uuid string) {}

func (ur *SqlUserRepository) Delete(uuid string) {}

func (ur *SqlUserRepository) FindByEmail(email string) {}

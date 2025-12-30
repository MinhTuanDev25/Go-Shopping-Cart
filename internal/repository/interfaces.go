package repository

import (
	"context"
	"go-shopping-cart/internal/db/sqlc"
)

type UserRepository interface {
	FindAll()
	Create(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	FindByUUID(uuid string)
	Update(uuid string)
	Delete(uuid string)
	FindByEmail(email string)
}

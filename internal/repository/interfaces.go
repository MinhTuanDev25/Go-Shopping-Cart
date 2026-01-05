package repository

import (
	"context"
	"go-shopping-cart/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindAll()
	Create(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	FindByUUID(uuid string)
	Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	Delete(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	SoftDelete(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	Restore(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	FindByEmail(email string)
}

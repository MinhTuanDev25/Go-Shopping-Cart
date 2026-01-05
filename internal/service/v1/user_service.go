package v1service

import (
	"database/sql"
	"errors"
	"go-shopping-cart/internal/db/sqlc"
	"go-shopping-cart/internal/repository"
	"go-shopping-cart/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (us *userService) GetAllUsers(search string, page, limit int) {
	// TODO: implement
}

func (us *userService) CreateUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	input.UserEmail = utils.NormalizeString(input.UserEmail)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, utils.WrapError(err, "Failed to hash password", utils.ErrCodeInternal)
	}
	input.UserPassword = string(hashedPassword)

	user, err := us.repo.Create(context, input)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.User{}, utils.NewError("Email already exists", utils.ErrCodeConflict)
		}

		return sqlc.User{}, utils.WrapError(err, "Failed to create user", utils.ErrCodeInternal)
	}

	return user, nil
}

func (us *userService) GetUserByUUID(uuid string) {
}

func (us *userService) UpdateUser(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	if input.UserPassword != nil && *input.UserPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			return sqlc.User{}, utils.WrapError(err, "Failed to hash password", utils.ErrCodeInternal)
		}
		hased := string(hashedPassword)
		input.UserPassword = &hased
	}

	userUpdated, err := us.repo.Update(context, input)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return sqlc.User{}, utils.NewError("Age must be greater than 0 and less than 150", utils.ErrCodeConflict)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "Failed to update user", utils.ErrCodeInternal)
	}

	return userUpdated, nil
}

func (us *userService) DeleteUser(ctx *gin.Context, uuid uuid.UUID) error {
	context := ctx.Request.Context()

	_, err := us.repo.Delete(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("user not found or not marked as delete for permenent removal", utils.ErrCodeNotFound)
		}

		return utils.WrapError(err, "failed to restore user", utils.ErrCodeInternal)
	}

	return nil
}

func (us *userService) SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	userSoftDelted, err := us.repo.SoftDelete(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}

		return sqlc.User{}, utils.WrapError(err, "failed to delete user", utils.ErrCodeInternal)
	}

	return userSoftDelted, nil

}

func (us *userService) RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	userRestored, err := us.repo.Restore(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found or not marked as delete for restore", utils.ErrCodeNotFound)
		}

		return sqlc.User{}, utils.WrapError(err, "failed to restore user", utils.ErrCodeInternal)
	}

	return userRestored, nil
}

package v1service

import (
	"database/sql"
	"errors"
	"go-shopping-cart/internal/db/sqlc"
	"go-shopping-cart/internal/repository"
	"go-shopping-cart/internal/utils"
	"strconv"

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

func (us *userService) GetAllUsers(ctx *gin.Context, search, orderBy, sort string, page, limit int32, deleted bool) ([]sqlc.User, int32, error) {
	context := ctx.Request.Context()

	if sort == "" {
		sort = "asc"
	}

	if orderBy == "" {
		orderBy = "user_created_at"
	}

	if page == 0 {
		page = 1
	}

	if limit <= 0 {
		limitStr := utils.GetEnv("LIMIT_RECORDS_PER_PAGE", "10")
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil || limitInt <= 0 {
			limit = 10
		}
		limit = int32(limitInt)
	}

	offset := (page - 1) * limit

	users, err := us.repo.GetAllV2(context, search, orderBy, sort, limit, offset, deleted)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError(err, "Failed to get all users", utils.ErrCodeInternal)
	}

	total, err := us.repo.CountUsers(context, search, deleted)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError(err, "Failed to count users", utils.ErrCodeInternal)
	}

	return users, int32(total), nil
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

func (us *userService) GetUserByUuid(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	user, err := us.repo.GetByUuid(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "failed to get user", utils.ErrCodeInternal)
	}

	return user, nil
}

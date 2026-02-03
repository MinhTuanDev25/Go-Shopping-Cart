package v1handler

import (
	v1dto "go-shopping-cart/internal/dto/v1"
	v1service "go-shopping-cart/internal/service/v1"
	"go-shopping-cart/internal/utils"
	"go-shopping-cart/internal/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service v1service.UserService
}

func NewUserHandler(service v1service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUsers(ctx *gin.Context) {
	var params v1dto.GetUsersParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	users, total, err := uh.service.GetAllUsers(ctx, params.Search, params.OrderBy, params.Sort, params.Page, params.Limit, false)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTOs := v1dto.MapUsersToDTO(users)
	pagination := utils.NewPaginationResponse(userDTOs, params.Page, params.Limit, total)

	utils.ResponseSuccess(ctx, http.StatusOK, "Users fetched successfully", pagination)
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var input v1dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapCreateInputToModel()

	newUser, err := uh.service.CreateUser(ctx, user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := v1dto.MapUserToDTO(newUser)

	utils.ResponseSuccess(ctx, http.StatusCreated, "User created successfully", userDTO)
}

func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	user, err := uh.service.GetUserByUuid(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := v1dto.MapUserToDTO(user)

	utils.ResponseSuccess(ctx, http.StatusOK, "User fetched successfully", userDTO)
}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	var input v1dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapUpdateInputToModel(userUuid)

	updatedUser, err := uh.service.UpdateUser(ctx, user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := v1dto.MapUserToDTO(updatedUser)

	utils.ResponseSuccess(ctx, http.StatusOK, "User updated successfully", userDTO)
}

func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	err = uh.service.DeleteUser(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}

func (uh *UserHandler) SoftDeleteUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
	}

	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	softDeleteUser, err := uh.service.SoftDeleteUser(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := v1dto.MapUserToDTO(softDeleteUser)

	utils.ResponseSuccess(ctx, http.StatusOK, "User deleted successfully", userDTO)
}

func (uh *UserHandler) RestoreUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	restoreUser, err := uh.service.RestoreUser(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := v1dto.MapUserToDTO(restoreUser)

	utils.ResponseSuccess(ctx, http.StatusOK, "User restored successfully", userDTO)
}

func (uh *UserHandler) GetSoftDeletedUsers(ctx *gin.Context) {
	var params v1dto.GetUsersParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	users, total, err := uh.service.GetAllUsers(ctx, params.Search, params.OrderBy, params.Sort, params.Page, params.Limit, true)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTOs := v1dto.MapUsersToDTO(users)
	pagination := utils.NewPaginationResponse(userDTOs, params.Page, params.Limit, total)

	utils.ResponseSuccess(ctx, http.StatusOK, "Soft deleted users fetched successfully", pagination)
}

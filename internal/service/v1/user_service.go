package v1service

import "user-management-api/internal/repository"

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (us *userService) GetAllUsers(search string, page, limit int) {
	// TODO: implement
}

func (us *userService) CreateUser() {
	// TODO: implement
}

func (us *userService) GetUserByUUID(uuid string) {
}

func (us *userService) UpdateUser(uuid string) {
	// TODO: implement
}

func (us *userService) DeleteUser(uuid string) {
	// TODO: implement
}

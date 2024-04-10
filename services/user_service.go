package services

import (
	"go-project/models"
	"go-project/repositories"
)

type CreateUserInput struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type UpdateUserInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type IUserService interface {
	Find() ([]models.User, error)
	FindOneById(id int) (models.User, error)
	FindOneByEmail(email string) (models.User, error)
	Create(input CreateUserInput) (models.User, error)
	Update(id int, input UpdateUserInput) (models.User, error)
	Delete(id int) error
}

type UserService struct {
	Repository repositories.IUserRepository
}

func NewUserService(repository repositories.IUserRepository) IUserService {
	return &UserService{Repository: repository}
}

func (service *UserService) Find() ([]models.User, error) {
	return service.Repository.Find()
}

func (service *UserService) FindOneById(id int) (models.User, error) {
	return service.Repository.FindOneById(id)
}

func (service *UserService) FindOneByEmail(email string) (models.User, error) {
	return service.Repository.FindOneByEmail(email)
}

func (service *UserService) Create(input CreateUserInput) (models.User, error) {
	user := models.User{FirstName: input.FirstName, LastName: input.LastName, Email: input.Email, Password: input.Password}
	return service.Repository.Create(user)
}

func (service *UserService) Update(id int, input UpdateUserInput) (models.User, error) {
	user, err := service.Repository.FindOneById(id)

	if err != nil {
		return models.User{}, err
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}

	return service.Repository.Update(user)
}

func (service *UserService) Delete(id int) error {
	return service.Repository.Delete(id)
}

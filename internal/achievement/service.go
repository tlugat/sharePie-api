package achievement

import (
	"sharePie-api/internal/models"
)

type CreateAchievementInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Points      int    `json:"points" binding:"required"`
	Condition   string `json:"condition" binding:"required"`
}

type UpdateAchievementInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Points      int    `json:"points"`
	Condition   string `json:"condition"`
}

type IAchievementService interface {
	Find() ([]models.Achievement, error)
	FindOne(id uint) (models.Achievement, error)
	Create(input CreateAchievementInput) (models.Achievement, error)
	Update(id uint, input UpdateAchievementInput) (models.Achievement, error)
	Delete(id uint) error
}

type Service struct {
	Repository IAchievementRepository
}

func NewService(
	repository IAchievementRepository,
) IAchievementService {
	return &Service{
		Repository: repository,
	}
}

func (service *Service) Find() ([]models.Achievement, error) {
	achievements, err := service.Repository.Find()
	if err != nil {
		return nil, err
	}

	return achievements, nil
}

func (service *Service) FindOne(id uint) (models.Achievement, error) {
	achievement, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Achievement{}, err
	}

	return achievement, nil
}

func (service *Service) Create(input CreateAchievementInput) (models.Achievement, error) {
	achievement := models.Achievement{
		Name:        input.Name,
		Description: input.Description,
		Points:      input.Points,
		Condition:   input.Condition,
	}

	return service.Repository.Create(achievement)
}

func (service *Service) Update(id uint, input UpdateAchievementInput) (models.Achievement, error) {
	achievement, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Achievement{}, err
	}

	if input.Name != "" {
		achievement.Name = input.Name
	}
	if input.Description != "" {
		achievement.Description = input.Description
	}
	if input.Points != 0 {
		achievement.Points = input.Points
	}
	if input.Condition != "" {
		achievement.Condition = input.Condition
	}

	return service.Repository.Update(achievement)
}

func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}
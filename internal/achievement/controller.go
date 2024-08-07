package achievement

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sharePie-api/internal/types"
	"strconv"
)

type Controller struct {
	achievementService types.IAchievementService
}

func NewController(service types.IAchievementService) *Controller {
	return &Controller{achievementService: service}
}

// FindAchievements retrieves all achievements.
// @Summary List all achievements
// @Description Retrieves a list of all achievements from the database
// @Tags Achievements
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Returns a list of achievements"
// @Failure 500 {object} map[string]interface{} "Returns an error if the request fails"
// @Router /achievements [get]
func (controller *Controller) FindAchievements(c *gin.Context) {
	achievements, err := controller.achievementService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve achievements: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": achievements})
}

// FindAchievement retrieves an achievement by ID.
// @Summary Get a single achievement
// @Description Retrieves an achievement by its ID from the database
// @Tags Achievements
// @Accept  json
// @Produce  json
// @Param id path int true "Achievement ID"
// @Success 200 {object} map[string.interface{} "Returns the specified achievement"
// @Failure 400 {object} map[string.interface{} "Returns an error if the achievement is not found"
// @Router /achievements/{id} [get]
func (controller *Controller) FindAchievement(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	achievement, err := controller.achievementService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Achievement with ID %d not found", id)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": achievement})
}

// CreateAchievement adds a new achievement.
// @Summary Add a new achievement
// @Description Adds a new achievement to the database, linked to the authenticated user
// @Tags Achievements
// @Accept  json
// @Produce  json
// @Param input body types.CreateAchievementInput true "Achievement creation data"
// @Success 200 {object} map[string.interface{} "Returns the newly created achievement"
// @Failure 400 {object} map[string.interface{} "Returns an error if the input is invalid or user authentication fails"
// @Router /achievements [post]
func (controller *Controller) CreateAchievement(c *gin.Context) {
	var input types.CreateAchievementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %v", err)})
		return
	}

	achievement, err := controller.achievementService.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create achievement: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": achievement})
}

// UpdateAchievement updates an existing achievement.
// @Summary Update an achievement
// @Description Updates an existing achievement with new data
// @Tags Achievements
// @Accept  json
// @Produce  json
// @Param id path int true "Achievement ID"
// @Param input body types.UpdateAchievementInput true "Achievement update data"
// @Success 200 {object} map[string.interface{} "Returns the updated achievement"
// @Failure 400 {object} map[string.interface{} "Returns an error if the input is invalid or the achievement does not exist"
// @Router /achievements/{id} [patch]
func (controller *Controller) UpdateAchievement(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input types.UpdateAchievementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %v", err)})
		return
	}
	achievement, err := controller.achievementService.Update(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update achievement with ID %d: %v", id, err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": achievement})
}

// DeleteAchievement removes an achievement.
// @Summary Delete an achievement
// @Description Deletes an achievement from the database
// @Tags Achievements
// @Accept  json
// @Produce  json
// @Param id path int true "Achievement ID"
// @Success 200 {object} map[string.interface{} "Confirms successful deletion"
// @Failure 400 {object} map[string.interface{} "Returns an error if the achievement cannot be deleted"
// @Router /achievements/{id} [delete]
func (controller *Controller) DeleteAchievement(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.achievementService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to delete achievement with ID %d: %v", id, err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}

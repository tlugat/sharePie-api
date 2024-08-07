package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sharePie-api/internal/auth"
	"sharePie-api/internal/models"
	"strconv"
)

func RefundIsUserPartOfEvent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := auth.GetUserFromContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		refundID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid refund ID"})
			return
		}

		var refund models.Refund
		if err := db.First(&refund, refundID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Refund not found"})
			return
		}

		var event models.Event
		if err := db.First(&event, refund.EventID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		result := db.Preload("Users").Where("id = ?", refund.EventID).First(&event)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Fetching user event error"})
			return
		}

		var users []models.User
		for _, user := range event.Users {
			fmt.Printf("user of event ==> %v\n", user.ID)
			users = append(users, user)
		}

		event.Users = users

		if !IsUserPartOfEvent(user, event) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User is not part of the event"})
			return
		}
		c.Next()
	}
}

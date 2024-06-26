package event

import (
	"errors"
	"net/http"
	"sharePie-api/internal/auth"
	"sharePie-api/internal/types"
	"sharePie-api/pkg/constants"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	eventService types.IEventService
}

func NewController(service types.IEventService) *Controller {
	return &Controller{eventService: service}
}

// FindEvents retrieves all events.
// @Summary List all events
// @Description Retrieves a list of all events from the database
// @Tags Events
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Returns a list of events"
// @Failure 500 {object} map[string]interface{} "Returns an error if the request fails"
// @Router /events [get]
func (controller *Controller) FindEvents(c *gin.Context) {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	if user.Role == constants.AdminRole {
		events, err := controller.eventService.Find()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": events})
		return
	}

	events, err := controller.eventService.FindByUser(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": events})

}

// FindEvent retrieves an event by ID.
// @Summary Get a single event
// @Description Retrieves an event by its ID from the database
// @Tags Events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]interface{} "Returns the specified event"
// @Failure 400 {object} map[string]interface{} "Returns an error if the event is not found"
// @Router /events/{id} [get]
func (controller *Controller) FindEvent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	event, err := controller.eventService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": event})
}

// CreateEvent adds a new event.
// @Summary Add a new event
// @Description Adds a new event to the database, linked to the authenticated user
// @Tags Events
// @Accept  json
// @Produce  json
// @Param input body services.CreateEventInput true "Event creation data"
// @Success 200 {object} map[string]interface{} "Returns the newly created event"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or user authentication fails"
// @Router /events [post]
func (controller *Controller) CreateEvent(c *gin.Context) {
	var input types.CreateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := auth.GetUserFromContext(c)

	if !ok {
		return
	}

	event, err := controller.eventService.Create(input, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": event})
}

// UpdateEvent updates an existing event.
// @Summary Update an event
// @Description Updates an existing event with new data
// @Tags Events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Param input body services.UpdateEventInput true "Event update data"
// @Success 200 {object} map[string]interface{} "Returns the updated event"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or the event does not exist"
// @Router /events/{id} [put]
func (controller *Controller) UpdateEvent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input types.UpdateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event, err := controller.eventService.Update(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": event})
}

// UpdateEventState updates an existing event.
// @Summary Update an event
// @Description Updates an existing event with new data
// @Tags Events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Param input body services.UpdateEventStateInput true "Event update data"
// @Success 200 {object} map[string]interface{} "Returns the updated event"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or the event does not exist"
// @Router /events/{id} [put]
func (controller *Controller) UpdateEventState(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input types.UpdateEventStateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event, err := controller.eventService.UpdateState(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": event})
}

// DeleteEvent removes an event.
// @Summary Delete an event
// @Description Deletes an event from the database
// @Tags Events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]interface{} "Confirms successful deletion"
// @Failure 400 {object} map[string]interface{} "Returns an error if the event cannot be deleted"
// @Router /events/{id} [delete]
func (controller *Controller) DeleteEvent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.eventService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete event"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}

// GetEventUsers retrieves all users for an event.
// @Summary Get event users
// @Description Retrieves all users for a specified event
// @Tags Events
// @Accept  json
// @Produce  json
// @Param eventId path int true "Event ID"
// @Success 200 {object} map[string]interface{} "Returns a list of users for the event"
// @Failure 400 {object} map[string]interface{} "Returns an error if the event does not exist"
// @Router /events/{id}/users [get]
func (controller *Controller) GetEventUsers(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	users, err := controller.eventService.GetUsers(uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetEventBalances retrieves a list of balances for an event.
// @Summary Get event balance list
// @Description Retrieves a summary of balances for a specified event
// @Tags Events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]interface{} "Returns a list of balances for the event"
// @Failure 400 {object} map[string]interface{} "Returns an error if the event does not exist"
// @Router /events/{id}/summary [get]
func (controller *Controller) GetEventBalances(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	event, err := controller.eventService.FindOne(uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event not found!"})
		return
	}

	balanceSummary, err := controller.eventService.GetBalances(event)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": balanceSummary})
}

// GetEventTransactions retrieves a list of transactions for an event.
// // @Summary Get event transactions
// // @Description Retrieves a list of transactions for a specified event
// // @Tags Events
// // @Accept  json
// // @Produce  json
// // @Param id path int true "Event ID"
// // @Success 200 {object} map[string]interface{} "Returns a list of transactions for the event"
// // @Failure 400 {object} map[string]interface{} "Returns an error if the event does not exist"
// // @Router /events/{id}/transactions [get]
func (controller *Controller) GetEventTransactions(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	event, err := controller.eventService.FindOne(uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event not found!"})
		return
	}

	transactions, err := controller.eventService.GetTransactions(event)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

func (controller *Controller) JoinEvent(c *gin.Context) {
	var input types.JoinEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := auth.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	event, err := controller.eventService.AddUser(input.Code, user)
	if err != nil {
		var conflictErr *types.ConflictError
		if errors.As(err, &conflictErr) {
			c.JSON(http.StatusConflict, gin.H{"error": conflictErr.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": event})
}

func (controller *Controller) GetExpenses(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	expense, err := controller.eventService.FindExpenses(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expense})
}

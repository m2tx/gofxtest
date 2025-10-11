package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m2tx/gofxtest/domain"
)

type eventRoute struct {
	RouteHandler
	eventService domain.EventService
}

func NewEventRoute(eventService domain.EventService) RouteHandler {
	return &eventRoute{
		eventService: eventService,
	}
}

func (r *eventRoute) Register(e *gin.Engine) {
	g := e.Group("/events")
	g.GET("", r.events)
	g.GET(":ID", r.eventByID)
	g.POST("", r.eventCreate)
	g.PUT(":ID", r.eventUpdate)
	g.DELETE(":ID", r.eventDelete)
}

// @Summary Create a new event
// @Description Create a new event with the provided details
// @Accept json
// @Produce json
// @Param event body domain.Event true "Event details"
// @Success 201 {object} domain.Event
// @Failure 400 {object} error "Invalid input"
// @Router /events [post]
func (r *eventRoute) eventCreate(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	var req domain.Event
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, err := r.eventService.Create(ctx, req)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary Get all events
// @Description Retrieve a list of all events
// @Accept json
// @Produce json
// @Success 200 {array} domain.Event
// @Router /events [get]
func (r *eventRoute) events(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	events, err := r.eventService.GetAll(ctx)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Get event by ID
// @Description Retrieve a single event by its ID
// @Param ID path string true "Event ID"
// @Accept json
// @Produce json
// @Success 200 {object} domain.Event
// @Failure 404 {object} error "Event not found"
// @Router /events/{ID} [get]
func (r *eventRoute) eventByID(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	id := c.Param("ID")

	event, err := r.eventService.Get(ctx, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, event)
}

// @Summary Update an existing event
// @Description Update the details of an existing event by its ID
// @Param ID path string true "Event ID"
// @Accept json
// @Produce json
// @Param event body domain.Event true "Updated event details"
// @Success 200 {object} domain.Event
// @Failure 400 {object} error "Invalid input"
// @Failure 404 {object} error "Event not found"
// @Router /events/{ID} [put]
func (r *eventRoute) eventUpdate(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	id := c.Param("ID")

	var req domain.Event
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	req.ID = id

	err := r.eventService.Update(ctx, req)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Delete an event
// @Description Delete an existing event by its ID
// @Param ID path string true "Event ID"
// @Accept json
// @Produce json
// @Success 204 "No Content"
// @Failure 404 {object} error "Event not found"
// @Router /events/{ID} [delete]
func (r *eventRoute) eventDelete(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	id := c.Param("ID")

	err := r.eventService.Delete(ctx, id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

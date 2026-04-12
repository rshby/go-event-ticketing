package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rshby/go-event-ticketing/internal/controller/response"
	"github.com/rshby/go-event-ticketing/internal/entity"
	"github.com/rshby/go-event-ticketing/tracing"
	"github.com/rshby/go-event-ticketing/utils/helper"
	"github.com/sirupsen/logrus"
)

type EventController struct {
	eventService entity.EventService
}

// NewEventController creates new instance of event controller
func NewEventController(eventService entity.EventService) *EventController {
	return &EventController{
		eventService: eventService,
	}
}

// CreateEvent creates new event
func (e *EventController) CreateEvent(c *gin.Context) {
	ctx, span := tracing.Start(c)
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
	})

	var request entity.CreateEventRequestDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error(err)
		response.ResponseError(c, err)
		return
	}

	// call method in service
	if err := e.eventService.CreateEvent(ctx, request); err != nil {
		logger.Error(err)
		response.ResponseError(c, err)
		return
	}

	response.NewResponse().
		WithMessage("success create event").
		ToResponseAPI(c, http.StatusCreated)
}

// GetListEvents retrieves event
func (e *EventController) GetListEvents(c *gin.Context) {
	ctx, span := tracing.Start(c)
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
	})

	// call method in service
	events, err := e.eventService.GetListEvents(ctx)
	if err != nil {
		logger.Error(err)
		response.ResponseError(c, err)
		return
	}

	logger.Infof("success retrieves events")
	response.NewResponse().
		WithMessage("success retrieves list events").
		WithData(&events).
		ToResponseAPI(c, http.StatusOK)
}

// GetEventByID retrieves event by id
func (e *EventController) GetEventByID(c *gin.Context) {
	ctx, span := tracing.Start(c)
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
	})

	// get if from url param
	id := c.Param("id")

	// call method in service
	events, err := e.eventService.GetEventByID(ctx, helper.ExpectedNumber[uint64](id))
	if err != nil {
		logger.Error(err)
		response.ResponseError(c, err)
		return
	}

	// success get event by id
	response.NewResponse().
		WithMessage("success retrieves event by id").
		WithData(&events).
		ToResponseAPI(c, http.StatusOK)
}

// DeleteEventByID deletes event by id
func (e *EventController) DeleteEventByID(c *gin.Context) {
	ctx, span := tracing.Start(c)
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
	})

	// get id from params
	id := c.Param("id")

	// call method in service
	if err := e.eventService.DeleteEventByID(ctx, helper.ExpectedNumber[uint64](id)); err != nil {
		logger.Error(err)
		response.ResponseError(c, err)
		return
	}

	// success delete
	response.NewResponse().
		WithMessage("success delete event").
		ToResponseAPI(c, http.StatusOK)
}

package service

import (
	"context"

	"github.com/rshby/go-event-ticketing/internal/entity"
	"github.com/rshby/go-event-ticketing/internal/errors"
	"github.com/rshby/go-event-ticketing/tracing"
	"github.com/rshby/go-event-ticketing/utils/helper"
	"github.com/sirupsen/logrus"
)

type eventService struct {
	eventRepository entity.EventRepository
}

// NewEventService creates new event repository
func NewEventService(eventRepository entity.EventRepository) entity.EventService {
	return &eventService{
		eventRepository: eventRepository,
	}
}

// CreateEvent creates new event and save to database
func (e *eventService) CreateEvent(ctx context.Context, request entity.CreateEventRequestDTO) error {
	ctx, span := tracing.Start(ctx)
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
		"request": helper.Dump(&request),
	})

	// call method in repository
	if err := e.eventRepository.Insert(ctx, request.ToEventEntity()); err != nil {
		logger.Error(err)
		return err
	}

	// success insert new event
	return nil
}

// GetListEvents retrieves list events from database
func (e *eventService) GetListEvents(ctx context.Context) ([]entity.GetListEventsResponseDTO, error) {
	ctx, span := tracing.Start(ctx)
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
	})

	// call method in repository
	events, err := e.eventRepository.GetListEvent(ctx)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// check if not found
	if len(events) == 0 {
		return nil, nil
	}

	var response = make([]entity.GetListEventsResponseDTO, 0, len(events))
	for _, event := range events {
		response = append(response, entity.GetListEventsResponseDTO{
			ID:          helper.ExpectedString(event.ID),
			Name:        event.Name,
			BannerImage: event.BannerImage,
			Location:    event.Location,
			StartTime:   helper.DateTimeToString(event.StartTime),
			EndTime:     helper.DateTimeToString(event.EndTime),
		})
	}

	logger.Info("success get list events")
	return response, nil
}

// GetEventByID retrieves event by id
func (e *eventService) GetEventByID(ctx context.Context, id uint64) (*entity.GetEventByIDResponseDTO, error) {
	ctx, span := tracing.Start(ctx)
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
		"id":      id,
	})

	// call method in repository
	event, err := e.eventRepository.GetByID(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// check if response in not found
	if event == nil {
		err = errors.ErrEventNotFound
		logger.Error(err)
		return nil, err
	}

	// mapping to response
	response := entity.GetEventByIDResponseDTO{
		ID:          helper.ExpectedString(event.ID),
		Name:        event.Name,
		BannerImage: event.BannerImage,
		Location:    event.Location,
		StartTime:   helper.DateTimeToString(event.StartTime),
		EndTime:     helper.DateTimeToString(event.EndTime),
		CreatedAt:   helper.DateTimeToString(event.CreatedAt),
		UpdatedAt:   helper.DateTimeToString(event.UpdatedAt),
		CreatedBy:   event.CreatedBy,
	}

	return &response, nil
}

// DeleteEventByID deletes event by id
func (e *eventService) DeleteEventByID(ctx context.Context, id uint64) error {
	ctx, span := tracing.Start(ctx)
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
		"id":      id,
	})

	// call method in repository
	if err := e.eventRepository.DeleteByID(ctx, id); err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("success delete event")
	return nil
}

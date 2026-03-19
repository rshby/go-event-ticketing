package repository

import (
	"context"
	"errors"

	"github.com/rshby/go-event-ticketing/cacher"
	"github.com/rshby/go-event-ticketing/config"
	"github.com/rshby/go-event-ticketing/internal/entity"
	"github.com/rshby/go-event-ticketing/tracing"
	"github.com/rshby/go-event-ticketing/utils/helper"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type eventRepository struct {
	db    *gorm.DB
	cache cacher.CacheManager
}

// NewEventRepository create new instance of event repository
func NewEventRepository(db *gorm.DB, cache cacher.CacheManager) entity.EventRepository {
	return &eventRepository{
		db:    db,
		cache: cache,
	}
}

// Insert inserts new record event to database
func (e *eventRepository) Insert(ctx context.Context, input entity.Event) error {
	ctx, span := tracing.Start(ctx)
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
		"input":   helper.Dump(&input),
	})

	if !input.IsIDExists() {
		// generate ID
		id, err := helper.GenerateSonyFlakeID()
		if err != nil {
			logger.Error(err)
			return err
		}

		// assign id
		input.ID = uint64(id)
	}

	// insert to database
	if err := e.db.WithContext(ctx).Model(&entity.Event{}).Create(&input).Error; err != nil {
		logger.Error(err)
		return err
	}

	// success insert new event
	return nil
}

// GetByID retrieves data event by id
func (e *eventRepository) GetByID(ctx context.Context, id uint64) (*entity.Event, error) {
	ctx, span := tracing.Start(ctx)
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
		"id":      id,
	})

	cacheKey := cacher.GetEventByIDCacheKey(id)
	if config.EnableCaching() {
		var event entity.Event
		if err := e.cache.Get(ctx, cacheKey, &event); err != nil {
			logger.Error(err)
		}

		if event.IsIDExists() {
			logger.Infof("returning data event with id [%d] from Redis Cache!", event.ID)
			return &event, nil
		}
	}

	// retrieves from database
	var event entity.Event
	if err := e.db.WithContext(ctx).Model(&entity.Event{}).Where("id=?", id).Limit(1).Take(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.Error(err)
		return nil, err
	}

	// set to redis cache
	if config.EnableCaching() {
		if err := e.cache.Store(ctx, cacheKey, cacher.NewItem(event)); err != nil {
			logger.Error(err)
		}
	}

	// success get from database
	return &event, nil
}

// GetListEvent retrieves list of events
func (e *eventRepository) GetListEvent(ctx context.Context) ([]entity.Event, error) {
	ctx, span := tracing.Start(ctx)
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": helper.DumpIncomingContext(ctx),
	})

	// retrieves from database
	var events []entity.Event
	if err := e.db.WithContext(ctx).Model(&entity.Event{}).Order("id DESC").Find(&events).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	// check if not found
	if len(events) == 0 {
		logger.Info(gorm.ErrRecordNotFound)
		return nil, nil
	}

	logrus.Info("success retrieves list of events")
	return events, nil
}

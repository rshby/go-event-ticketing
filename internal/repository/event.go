package repository

import (
	"context"

	"github.com/rshby/go-event-ticketing/cacher"
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

	}
}

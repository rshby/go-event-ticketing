package entity

import (
	"context"
	"time"

	"github.com/rshby/go-event-ticketing/utils/helper"
)

type Event struct {
	ID          uint64    `gorm:"primaryKey;column:id"`
	Name        string    `gorm:"type:varchar(255);column:name;not null"`
	BannerImage string    `gorm:"type:varchar(255);column:banner_image"`
	Location    string    `gorm:"type:varchar(255);column:location"`
	StartTime   time.Time `gorm:"column:start_time"`
	EndTime     time.Time `gorm:"column:end_time"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedBy   string    `gorm:"type:varchar(255);column:created_by"`
}

func (e *Event) TableName() string {
	return "events"
}

func (e *Event) IsIDExists() bool {
	return e.ID > 0
}

type EventRepository interface {
	Insert(ctx context.Context, input Event) error
	GetByID(ctx context.Context, id uint64) (*Event, error)
	GetListEvent(ctx context.Context) ([]Event, error)
}

type EventService interface {
	CreateEvent(ctx context.Context, request CreateEventRequestDTO) error
	GetListEvents(ctx context.Context) ([]GetListEventsResponseDTO, error)
	GetEventByID(ctx context.Context, id uint64) (*GetEventByIDResponseDTO, error)
}

type CreateEventRequestDTO struct {
	Name        string `json:"name"`
	BannerImage string `json:"banner_image"`
	Location    string `json:"location"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

func (c *CreateEventRequestDTO) ToEventEntity() Event {
	now := helper.GetNowTime()
	id, _ := helper.GenerateSonyFlakeID()

	event := Event{
		ID:          uint64(id),
		Name:        c.Name,
		BannerImage: c.BannerImage,
		Location:    c.Location,
		StartTime:   helper.StringToDateTime(c.StartTime),
		EndTime:     helper.StringToDateTime(c.EndTime),
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy:   "admin",
	}

	return event
}

type GetListEventsResponseDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BannerImage string `json:"bannerImage"`
	Location    string `json:"location"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
}

type GetEventByIDResponseDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BannerImage string `json:"bannerImage"`
	Location    string `json:"location"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	CreatedBy   string `json:"createdBy"`
}

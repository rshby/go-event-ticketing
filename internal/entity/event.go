package entity

import "time"

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
}

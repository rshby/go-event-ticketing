package entity

import "time"

type Ticket struct {
	ID             uint64    `gorm:"primaryKey;column:id"`
	EventID        uint64    `gorm:"column:event_id;index"`
	Name           string    `gorm:"type:varchar(255);column:name;not null"`
	BannerImage    string    `gorm:"type:varchar(255);column:banner_image"`
	MaxQuota       int       `gorm:"type:integer;column:max_quota;not null"`
	AvailableQuota int       `gorm:"type:integer;column:available_quota;not null"`
	StartTime      time.Time `gorm:"column:start_time"`
	EndTime        time.Time `gorm:"column:end_time"`
	Price          float64   `gorm:"type:decimal(12,2);column:price;not null"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (t *Ticket) TableName() string {
	return "tickets"
}

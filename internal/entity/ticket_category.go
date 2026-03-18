package entity

import "time"

type TicketCategory struct {
	ID        uint64    `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"type:varchar(255);column:name;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (t *TicketCategory) TableName() string {
	return "ticket_categories"
}

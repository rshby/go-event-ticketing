package entity

import "time"

type Order struct {
	ID               uint64     `gorm:"primaryKey;autoIncrement;column:id"`
	TicketID         uint64     `gorm:"column:ticket_id;index"`
	TicketCategoryID uint64     `gorm:"column:ticket_category_id;index"`
	OrderID          string     `gorm:"type:varchar(100);column:order_id;uniqueIndex;not null"`
	Email            string     `gorm:"type:varchar(255);column:email"`
	Phone            string     `gorm:"type:varchar(50);column:phone"`
	Name             string     `gorm:"type:varchar(255);column:name"`
	Gender           string     `gorm:"type:varchar(10);column:gender"`
	BibNumber        string     `gorm:"type:varchar(50);column:bib_number"`
	BloodType        string     `gorm:"type:varchar(10);column:blood_type"`
	EmergencyCall    string     `gorm:"type:varchar(50);column:emergency_call"`
	PaymentStatus    string     `gorm:"type:varchar(50);column:payment_status;not null"`
	JerseySize       string     `gorm:"type:varchar(10);column:jersey_size"`
	PriceAmount      float64    `gorm:"type:decimal(12,2);column:price_amount;not null"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	ExpiredAt        *time.Time `gorm:"column:expired_at"`
}

func (o *Order) TableName() string {
	return "orders"
}

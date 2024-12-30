package models

import "time"

type Order struct {
	Id               int        `json:"id" form:"id" gorm:"primaryKey"`
	UserID           int        `json:"userId" gorm:"not null"`
	TotalPrice       float64    `json:"total_price"`
	Status           string     `json:"status"`
	User             User       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	PaymentCreatedAt *time.Time `json:"paymentCreatedAt" form:"updatedAt" gorm:"autoUpdateTime:false"`
	PaymentDeadline  time.Time
	CreatedAt        time.Time     `json:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt"`
	OrderDetails     []OrderDetail `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type OrderDetail struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	OrderID   int       `json:"orderId" gorm:"not null"`
	TicketID  int       `json:"ticketId" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Subtotal  float64   `json:"subtotal" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time `json:"updatedAt"`
	Order     Order     `gorm:"constraint:OnDelete:CASCADE" json:"-"` // Relasi ke Order
	Ticket    Ticket    `gorm:"constraint:OnDelete:CASCADE"`          // Relasi ke Ticket
}

type PayloadOrder struct {
	UserID      int                  `json:"userID"`
	OrderDetail []PayloadOrderDetail `json:"orderDetails"`
}

type PayloadOrderDetail struct {
	TicketID int `json:"ticketId"`
	Qty      int `json:"qty"`
}

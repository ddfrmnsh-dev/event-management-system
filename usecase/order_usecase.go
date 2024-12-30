package usecase

import (
	"event-management-system/models"
	"event-management-system/repository"
	"fmt"
	"log"
	"time"
)

type OrderUseCase interface {
	CreateOrder(order models.PayloadOrder) (models.Order, error)
	ExpireUnpaidOrders() error
}

type orderUseCaseImpl struct {
	orderRepository repository.OrderRepository
}

func NewOrderUseCase(orderRepository repository.OrderRepository) OrderUseCase {
	return &orderUseCaseImpl{orderRepository: orderRepository}
}

func (oc *orderUseCaseImpl) CreateOrder(order models.PayloadOrder) (models.Order, error) {
	tx := oc.orderRepository.BeginTransaction()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := models.User{}

	if err := tx.First(&user, order.UserID).Error; err != nil {
		tx.Rollback()
		return models.Order{}, fmt.Errorf("invalid user ID")
	}

	var totalPrice float64
	var orderTickets []models.OrderDetail

	for _, ticket := range order.OrderDetail {
		var t models.Ticket
		log.Println("CHECK QTY", ticket.Qty)
		log.Println("CHECK ID TICKET", ticket.TicketID)
		if err := tx.First(&t, ticket.TicketID).Error; err != nil {
			tx.Rollback()
			return models.Order{}, fmt.Errorf("invalid ticket ID: %d", ticket.TicketID)
		}

		if t.Quota < ticket.Qty {
			tx.Rollback()
			return models.Order{}, fmt.Errorf("insufficient quota for ticket ID: %d", ticket.TicketID)
		}

		subtotal := float64(ticket.Qty) * float64(t.Price)
		totalPrice += subtotal

		t.Quota -= ticket.Qty
		if err := tx.Save(&t).Error; err != nil {
			tx.Rollback()
			return models.Order{}, fmt.Errorf("failed to update ticket quota")
		}

		orderTickets = append(orderTickets, models.OrderDetail{
			TicketID: ticket.TicketID,
			Quantity: ticket.Qty,
			Subtotal: subtotal,
		})
	}

	orders := models.Order{
		UserID:          order.UserID,
		TotalPrice:      totalPrice,
		Status:          "Pending",
		PaymentDeadline: time.Now().Add(2 * time.Minute),
	}

	if err := tx.Create(&orders).Error; err != nil {
		tx.Rollback()
		log.Println("failed create order", err.Error())
		return models.Order{}, fmt.Errorf("failed to create order")
	}

	for i := range orderTickets {
		orderTickets[i].OrderID = orders.Id
		if err := tx.Create(&orderTickets[i]).Error; err != nil {
			tx.Rollback()
			return models.Order{}, fmt.Errorf("failed to create order tickets")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return models.Order{}, fmt.Errorf("failed to commit transaction")
	}

	return orders, nil
}

func (oc *orderUseCaseImpl) ExpireUnpaidOrders() error {
	now := time.Now()
	tx := oc.orderRepository.BeginTransaction()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var expiredOrders []models.Order

	if err := tx.Where("status = ? AND payment_deadline < ?", "Pending", now).Find(&expiredOrders).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to find unpaid orders: %w", err)
	}

	for _, order := range expiredOrders {
		var orderDetails []models.OrderDetail

		if err := tx.Where("order_id = ?", order.Id).Find(&orderDetails).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to fetch order details: %w", err)
		}

		for _, detail := range orderDetails {
			var ticket models.Ticket

			// Ambil data tiket terkait
			if err := tx.First(&ticket, detail.TicketID).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to fetch ticket: %w", err)
			}

			// Kembalikan kuota tiket
			ticket.Quota += detail.Quantity
			if err := tx.Save(&ticket).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update ticket quota: %w", err)
			}
		}

		// Ubah status order menjadi "Expired"
		order.Status = "Expired"
		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update order status: %w", err)
		}
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

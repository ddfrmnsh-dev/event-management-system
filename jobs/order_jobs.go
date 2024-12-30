package jobs

import (
	"event-management-system/usecase"
	"fmt"
)

type SchedulerOrderJobs interface {
	CheckStatusPayment()
}

type schedulerOrderJobs struct {
	orderUseCase usecase.OrderUseCase
}

func NewSchedulerOrderJobs(orderUseCase usecase.OrderUseCase) SchedulerOrderJobs {
	return &schedulerOrderJobs{
		orderUseCase: orderUseCase,
	}
}

func (s *schedulerOrderJobs) CheckStatusPayment() {
	err := s.orderUseCase.ExpireUnpaidOrders()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

}

package jobs

import (
	"event-management-system/usecase"
	"fmt"
	"log"
	"time"
)

type SchedulerJobs interface {
	SendEmailActivation()
}

type schedulerJobs struct {
	userUseCase usecase.UserUseCase
}

func NewSchedulerJobs(userUseCase usecase.UserUseCase) SchedulerJobs {
	return &schedulerJobs{
		userUseCase: userUseCase,
	}
}

func (s *schedulerJobs) SendEmailActivation() {
	userCheck, err := s.userUseCase.FinByParams("is_active", false)
	if err != nil {
		fmt.Println("Error fetching users:", err)
		return
	}

	for _, user := range userCheck {
		// fmt.Printf("Mengirim notifikasi ke pengguna: %s pada jam %s\n", user.Username, time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local))
		fmt.Printf("Mengirim notifikasi ke pengguna: %s pada jam %s\n", user.Username, time.TimeOnly)
		log.Printf("Mengirim notifikasi ke pengguna: %s pada jam %s\n", user.Username, time.TimeOnly)
	}

}

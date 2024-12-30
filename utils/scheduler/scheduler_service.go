package scheduler

import (
	"event-management-system/config"
	"event-management-system/jobs"
	"fmt"
)

type SchedulerService interface {
	SendEmailActivation() error
	CheckPaymentOrder() error
}
type schedulerService struct {
	cfg                config.SchedulerConfig
	jobsScheduler      jobs.SchedulerJobs
	jobsOrderScheduler jobs.SchedulerOrderJobs
}

func NewSchedulerService(cfg config.SchedulerConfig, jobsScheduler jobs.SchedulerJobs, jobsOrderScheduler jobs.SchedulerOrderJobs) SchedulerService {
	return &schedulerService{
		cfg:                cfg,
		jobsScheduler:      jobsScheduler,
		jobsOrderScheduler: jobsOrderScheduler,
	}
}

// func (s *schedulerService) TestingScheduler() error {
// 	fmt.Println("Scheduler jalan yaa!!")
// 	cronExpression := "*/5 * * * *"

// 	_, err := s.cfg.Cron.AddFunc(cronExpression, sendNotifications)
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("Job scheduled every hours")
// 	// s.cfg.Cron.Start()
// 	return nil
// }

func (s *schedulerService) SendEmailActivation() error {
	// defer s.cfg.Cron.Stop()
	fmt.Println("SCHEDULER EMAIL")

	cronExpression := "*/5 * * * *"
	_, err := s.cfg.Cron.AddFunc(cronExpression, s.jobsScheduler.SendEmailActivation)
	if err != nil {
		return err
	}
	s.cfg.Cron.Start()
	return nil
}
func (s *schedulerService) CheckPaymentOrder() error {
	fmt.Println("SCHEDULER CHECK PAYMENT")

	cronExpression := "*/1 * * * *"
	_, err := s.cfg.Cron.AddFunc(cronExpression, s.jobsOrderScheduler.CheckStatusPayment)
	if err != nil {
		return err
	}
	s.cfg.Cron.Start()
	return nil
}

package main

import (
	"event-management-system/config"
	"event-management-system/controller"
	"event-management-system/jobs"
	"event-management-system/middleware"
	"event-management-system/models"
	"event-management-system/repository"
	"event-management-system/usecase"
	"event-management-system/utils/scheduler"
	"event-management-system/utils/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Server struct {
	userUC           usecase.UserUseCase
	authUC           usecase.AuthenticationUseCase
	eventUC          usecase.EventUseCase
	ticketUC         usecase.TicketUseCase
	orderUC          usecase.OrderUseCase
	jwtService       service.JwtService
	schedulerService scheduler.SchedulerService
	engine           *gin.Engine
	host             string
}

func (s *Server) initRoute() {
	rgAuth := s.engine.Group("/api/auth")
	controller.NewAuthController(s.authUC, rgAuth).Route()

	rgV1 := s.engine.Group("/api/v1")
	authMiddleware := middleware.NewAuthMiddleware(s.jwtService)

	controller.NewUserController(s.userUC, rgV1, authMiddleware).Route()
	controller.NewEventController(s.eventUC, rgV1, authMiddleware).Route()
	controller.NewTicketController(s.ticketUC, rgV1, authMiddleware).Route()
	controller.NewOrderController(s.orderUC, rgV1, authMiddleware).Route()
}

func (s *Server) initMigration() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Ticket{},
		&models.Order{},
		&models.OrderDetail{},
	)

	if err != nil {
		log.Fatal("Failed to migrate")
	}

	fmt.Println("Migrated Successfully")
}

func (s *Server) initScheduler() {
	if err := s.schedulerService.SendEmailActivation(); err != nil {
		log.Fatalf("Failed to initialize scheduler: %v", err)
	}

	if err := s.schedulerService.CheckPaymentOrder(); err != nil {
		log.Fatalf("Failed to initialize scheduler: %v", err)
	}
}

func (s *Server) Run() {
	s.initRoute()
	s.initMigration()
	s.initScheduler()
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("server not running on host %s, because error %v", s.host, err))
	}
}

func NewServer() *Server {
	var err error

	cfg, _ := config.NewConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("connection error")
	}

	userRepo := repository.NewUserRepository(DB)
	eventRepo := repository.NewEventRepository(DB)
	ticketRepo := repository.NewTicketRepository(DB)
	orderRepo := repository.NewOrderRepository(DB)

	userUseCase := usecase.NewUserUseCase(userRepo)
	eventUseCase := usecase.NewEventUseCase(eventRepo, userRepo)
	ticketUseCase := usecase.NewTicketUseCase(ticketRepo)
	orderUseCase := usecase.NewOrderUseCase(orderRepo)

	jwtService := service.NewJwtService(cfg.TokenConfig)
	schedulerJobs := jobs.NewSchedulerJobs(userUseCase)
	schedulerOrderJobs := jobs.NewSchedulerOrderJobs(orderUseCase)
	schedulerService := scheduler.NewSchedulerService(cfg.SchedulerConfig, schedulerJobs, schedulerOrderJobs)
	authUseCase := usecase.NewAuthenticationUseCase(userUseCase, jwtService)
	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	return &Server{
		userUC:           userUseCase,
		authUC:           authUseCase,
		eventUC:          eventUseCase,
		ticketUC:         ticketUseCase,
		orderUC:          orderUseCase,
		engine:           engine,
		jwtService:       jwtService,
		schedulerService: schedulerService,
		host:             host,
	}
}

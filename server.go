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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// Swagger UI files
	_ "event-management-system/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // Gin Swagger middleware
)

var DB *gorm.DB

// var (
// 	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
// 		Name: "myapp_processed_ops_total",
// 		Help: "The total number of processed events",
// 	})
// )

// @title My API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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
	s.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	s.engine.Use(middleware.PrometheusMiddleware())
	rgAuth := s.engine.Group("/api/auth")
	controller.NewAuthController(s.authUC, rgAuth).Route()

	rgV1 := s.engine.Group("/api/v1")
	authMiddleware := middleware.NewAuthMiddleware(s.jwtService)
	controller.NewUserController(s.userUC, rgV1, authMiddleware).Route()
	controller.NewEventController(s.eventUC, rgV1, authMiddleware).Route()
	controller.NewTicketController(s.ticketUC, rgV1, authMiddleware).Route()
	controller.NewOrderController(s.orderUC, rgV1, authMiddleware).Route()

	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.engine.LoadHTMLGlob("templates/*")

	s.engine.GET("/visualization", func(c *gin.Context) {
		c.HTML(http.StatusOK, "visualization.html", nil)
	})

}

// func (s *Server) initRecordMetrics() {
// 	go func() {
// 		for {
// 			opsProcessed.Inc()
// 			time.Sleep(2 * time.Second)
// 		}
// 	}()
// }

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
	// if err := s.schedulerService.SendEmailActivation(); err != nil {
	// 	log.Fatalf("Failed to initialize scheduler: %v", err)
	// }

	if err := s.schedulerService.CheckPaymentOrder(); err != nil {
		log.Fatalf("Failed to initialize scheduler: %v", err)
	}
}

func (s *Server) Run() {
	s.initRoute()
	s.initMigration()
	s.initScheduler()
	// s.initRecordMetrics()
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

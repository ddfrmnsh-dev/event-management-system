package main

import (
	"event-management-system/config"
	"event-management-system/controller"
	"event-management-system/middleware"
	"event-management-system/models"
	"event-management-system/repository"
	"event-management-system/usecase"
	"event-management-system/utils/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Server struct {
	userUC     usecase.UserUseCase
	authUC     usecase.AuthenticationUseCase
	jwtService service.JwtService
	engine     *gin.Engine
	host       string
}

func (s *Server) initRoute() {
	rgAuth := s.engine.Group("/api/auth")
	controller.NewAuthController(s.authUC, rgAuth).Route()

	rgV1 := s.engine.Group("/api/v1")
	authMiddleware := middleware.NewAuthMiddleware(s.jwtService)
	controller.NewUserController(s.userUC, rgV1, authMiddleware).Route()
}

func (s *Server) initMigration() {
	err := DB.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		log.Fatal("Failed to migrate")
	}

	fmt.Println("Migrated Successfully")
}

func (s *Server) Run() {
	s.initRoute()
	s.initMigration()
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
	userUseCase := usecase.NewUserUseCase(userRepo)
	jwtService := service.NewJwtService(cfg.TokenConfig)
	authUseCase := usecase.NewAuthenticationUseCase(userUseCase, jwtService)
	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	return &Server{
		userUC:     userUseCase,
		authUC:     authUseCase,
		engine:     engine,
		jwtService: jwtService,
		host:       host,
	}
}

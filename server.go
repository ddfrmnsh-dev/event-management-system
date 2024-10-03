package main

import (
	"event-management-system/config"
	"event-management-system/models"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Server struct {
	engine *gin.Engine
	host   string
}

func (s *Server) initRoute() {
	// rg := s.engine.Group("/api/v1")
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

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	return &Server{
		engine: engine,
		host:   host,
	}
}

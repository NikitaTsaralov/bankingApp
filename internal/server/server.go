package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/pkg/rabbitmq"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	echo   *echo.Echo
	cfg    *config.Config
	broker *rabbitmq.RabbitMQClient
	db     *gorm.DB // or it will be inside service
	logger *log.Logger
}

func Init(cfg *config.Config, db *gorm.DB, broker *rabbitmq.RabbitMQClient, logger *log.Logger) *Server {
	return &Server{
		echo:   echo.New(),
		cfg:    cfg,
		db:     db,
		broker: broker,
		logger: logger,
	}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	s.echo.POST("/putMoney", s.putMoney)
	s.echo.POST("/getMoney", s.getMoney)

	s.echo.POST("/login", s.Login)
	s.echo.PUT("/register", s.Register)

	go func() {
		log.Printf("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			log.Fatal("Error starting Server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	log.Println("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}

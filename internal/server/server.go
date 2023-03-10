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
	"github.com/NikitaTsaralov/bankingApp/pkg/rabbitmq/reconnect"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	echo             *echo.Echo
	cfg              *config.Config
	rabbitConnection *reconnect.Connection
	db               *gorm.DB // or it will be inside service
	logger           *log.Logger
}

func Init(cfg *config.Config, db *gorm.DB, rabbitConnection *reconnect.Connection, logger *log.Logger) *Server {
	return &Server{
		echo:             echo.New(),
		cfg:              cfg,
		db:               db,
		rabbitConnection: rabbitConnection,
		logger:           logger,
	}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		log.Printf("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			log.Fatal("Error starting Server: ", err)
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	log.Println("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}

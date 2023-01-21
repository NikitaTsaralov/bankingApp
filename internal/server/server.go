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
	db     *gorm.DB // or it will be inside service
	logger *log.Logger
}

func Init(cfg *config.Config, db *gorm.DB) *Server {
	return &Server{
		cfg:    cfg,
		db:     db,
		logger: &log.Logger{},
	}
}

func (server *Server) Run() error {
	httpServer := &http.Server{
		Addr:           server.cfg.Server.Port,
		ReadTimeout:    time.Second * server.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * server.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		server.logger.Printf("Server is listening on PORT: %s", server.cfg.Server.Port)
		if err := server.echo.StartServer(httpServer); err != nil {
			server.logger.Fatal("Error starting Server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	server.logger.Println("Server Exited Properly")
	return server.echo.Server.Shutdown(ctx)
}

package server

import (
	"github.com/NikitaTsaralov/bankingApp/internal/middleware"
	userHttp "github.com/NikitaTsaralov/bankingApp/internal/users/delivery/http"
	userRepo "github.com/NikitaTsaralov/bankingApp/internal/users/repository"
	userUseCase "github.com/NikitaTsaralov/bankingApp/internal/users/usecase"

	// transactionHttp "github.com/NikitaTsaralov/bankingApp/internal/transactions/delivery/http"

	// transactionRepo "github.com/NikitaTsaralov/bankingApp/internal/transactions/repository"
	// transactionUseCase "github.com/NikitaTsaralov/bankingApp/internal/transactions/usecase"
	"github.com/labstack/echo/v4"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// publisher, err := rabbitmq.InitTransactionPublisher(s.cfg, s.logger)
	// if err != nil {
	// 	return err
	// }
	// defer publisher.Close()

	// register repos
	uRepo := userRepo.Init(s.db)
	// tRepo := transactionRepo.Init(s.db)

	// register usecases
	userUseCase := userUseCase.Init(s.cfg, uRepo, s.logger)
	// transactionUseCase := transactionUseCase.Init(s.cfg, tRepo, uRepo, publisher, s.logger)

	// register handlers
	userHandlers := userHttp.Init(s.cfg, userUseCase, s.logger)
	// transactionHandlers := transactionHttp.Init(s.cfg, transactionUseCase, s.logger)

	// register middlewares
	mw := middleware.Init(s.cfg, userUseCase, s.logger)

	v1 := e.Group("/api/v1")
	userGroup := v1.Group("/user")
	// transactionGroup := v1.Group("/transaction")

	// register routes
	userHttp.MapAuthRoutes(userGroup, userHandlers, mw)
	// transactionHttp.MapTransactionRoutes(transactionGroup, transactionHandlers, mw)

	return nil
}

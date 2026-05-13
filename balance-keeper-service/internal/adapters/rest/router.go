package rest

import (
	"fmt"

	"github.com/PlopyBlopy/balance-keeper-service/internal/adapters/postgres"
	"github.com/PlopyBlopy/balance-keeper-service/internal/adapters/rest/handlers"
	"github.com/PlopyBlopy/balance-keeper-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(v int, pool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	accountRepo := postgres.NewAccountRepository(pool)
	outboxRepo := postgres.NewOutboxRepository(pool)
	txManager := postgres.NewTxManager(pool)

	addAccountUsecase := usecase.AddAccount(txManager, accountRepo, outboxRepo)

	vgroup := r.Group(fmt.Sprintf("/v%d", v))
	vgroup.POST("/", handlers.AddAccount(addAccountUsecase))

	return r
}

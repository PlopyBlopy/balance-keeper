package integration

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/PlopyBlopy/balance-keeper-service/internal/adapters/postgres"
	"github.com/PlopyBlopy/balance-keeper-service/internal/domain"
	"github.com/PlopyBlopy/balance-keeper-service/internal/shared"
	"github.com/PlopyBlopy/balance-keeper-service/tests/testdata"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// В этом тесте используется Testcontainer где в начале каждого теста вызывается testSuite.SetupTestPg - он закрывает все подключения, делает Rollback до начального состояния
func TestAccountRepository(t *testing.T) {
	pgContainerCtx := context.Background()
	// config
	c, err := testdata.NewTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	// container
	pgContainer, err := testdata.NewPostgresTestcontainer(pgContainerCtx, *c)
	if err != nil {
		t.Fatal(err)
	}

	// connStr
	connStr, err := pgContainer.ConnectionString(pgContainerCtx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// project root from sys root
	root, err := shared.FindProjectRoot()
	if err != nil {
		t.Fatal(err)
	}

	// path to migrate folder
	migrationsPath := filepath.Join(root, "migrations")

	// initialize migration on db
	err = testdata.InitMigrate(migrationsPath, connStr)
	if err != nil {
		t.Fatal(err)
	}

	// wrap db conns in TestSuite
	testSuite, err := testdata.NewTestSuite(pgContainerCtx, pgContainer, *c)
	if err != nil {
		t.Fatal(err)
	}
	// remove testcontainer from docker
	t.Cleanup(func() {
		if err := testSuite.PgContainer.Terminate(pgContainerCtx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	t.Run("AddAccount", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		err := testSuite.SetupTestPg(ctx)
		require.NoError(t, err)

		assert := assert.New(t)
		require := require.New(t)

		accountRep := postgres.NewAccountRepository(testSuite.Db)
		outboxRep := postgres.NewOutboxRepository(testSuite.Db)

		account := domain.NewAccount(uuid.New())
		msg, err := domain.NewAccountCreatedEvent(
			domain.AccountCreatedEvent{
				Id:             account.Id,
				InitialBalance: domain.BalanceParseInt64(account.Balance),
				CreatedAt:      account.CreatedAt,
			},
		)

		// Act
		tx, err := testSuite.Db.Begin(ctx)
		t.Cleanup(func() {
			err := tx.Rollback(ctx)
			assert.ErrorIs(err, pgx.ErrTxClosed)
		})

		err = accountRep.AddAccountTx(tx, account, ctx)
		require.NoError(err)

		require.NoError(err)

		err = outboxRep.InsertTx(tx, msg, ctx)
		require.NoError(err)

		tx.Commit(ctx)

		actAccount, err := accountRep.GetAccount(account.Id, ctx)
		require.NoError(err)

		actMsg, err := outboxRep.Get(msg.Id, ctx)
		require.NoError(err)

		// Assert
		assert.Equal(account.Id, actAccount.Id)
		assert.Equal(msg.Id, actMsg.Id)
	})
}

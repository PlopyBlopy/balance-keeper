package testdata

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Creates postgres container
func NewPostgresTestcontainer(ctx context.Context, c TestConfig) (*postgres.PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx, c.Image,
		postgres.WithDatabase(c.Database),
		postgres.WithUsername(c.Username),
		postgres.WithPassword(c.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		return nil, err
	}

	return pgContainer, nil
}

type TestSuite struct {
	PgContainer  *postgres.PostgresContainer
	Db           *pgxpool.Pool
	c            TestConfig
	snapshotName string
}

func NewTestSuite(ctx context.Context, pgContainer *postgres.PostgresContainer, c TestConfig) (*TestSuite, error) {
	snapshotName := c.ShapshotName

	err := pgContainer.Snapshot(ctx, postgres.WithSnapshotName(snapshotName))
	if err != nil {
		return nil, err
	}

	pool, err := newPostgresConnection(ctx, pgContainer, c)
	if err != nil {
		return nil, err
	}

	return &TestSuite{
		PgContainer:  pgContainer,
		Db:           pool,
		c:            c,
		snapshotName: snapshotName,
	}, nil
}

// Completes active connections, applies base snapshot, restoring the database to its original form.
func (ts *TestSuite) SetupTestPg(ctx context.Context) error {
	closeConns(ts.Db)

	err := ts.PgContainer.Restore(ctx, postgres.WithSnapshotName(ts.snapshotName))
	if err != nil {
		return err
	}

	pool, err := newPostgresConnection(ctx, ts.PgContainer, ts.c)
	if err != nil {
		return err
	}

	ts.Db = pool

	return nil
}

func newPostgresConnection(ctx context.Context, pgContainer *postgres.PostgresContainer, c TestConfig) (*pgxpool.Pool, error) {
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	connConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	connConfig.MaxConns = c.MaxConns
	connConfig.MinConns = c.MinConns

	pool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func closeConns(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}

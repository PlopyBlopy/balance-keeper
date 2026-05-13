package postgres

import (
	"context"
	"errors"

	"github.com/PlopyBlopy/balance-keeper-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRepository struct {
	db
}

func NewOutboxRepository(pool *pgxpool.Pool) *OutboxRepository {
	return &OutboxRepository{
		db: db{pool: pool},
	}
}

func (r *OutboxRepository) InsertTx(tx pgx.Tx, msg domain.OutboxMessage, ctx context.Context) error {
	ct, err := tx.Exec(ctx, "INSERT INTO outbox (id, aggregate_id, event_type, payload, status, created_at) VALUES($1,$2,$3,$4,$5,$6)", msg.Id, msg.AggregateId, msg.EventType, msg.Payload, msg.Status, msg.CreateAt)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotAdded
	}

	return nil
}

func (r *OutboxRepository) Get(id uuid.UUID, ctx context.Context) (domain.OutboxMessage, error) {
	msg := domain.OutboxMessage{}

	err := r.pool.QueryRow(ctx, "SELECT * FROM outbox WHERE id = $1", id).Scan(&msg.Id, &msg.AggregateId, &msg.EventType, &msg.Payload, &msg.Status, &msg.CreateAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return msg, domain.ErrNotFound
		}
		return msg, err
	}

	return msg, nil
}

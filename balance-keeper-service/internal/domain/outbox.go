package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OutboxStatus string

var (
	OutboxStatusPending OutboxStatus = "pending" // в обработке
	OutboxStatusSent    OutboxStatus = "sent"    // отправлен
	OutboxStatusFailed  OutboxStatus = "failed"  // ошибка

)

type OutboxEvent string

var (
	OutboxEventAccountCreated      OutboxEvent = "account_created" // создан новый аккаунт
	OutboxEventBalanceDeposited    OutboxEvent = ""                // Зачисление на счет
	OutboxEventBalanceWithdraw     OutboxEvent = ""                // Списание с счета
	OutboxEventTransferCompleted   OutboxEvent = ""                // перевод выполнен
	OutboxEventTransferUncompleted OutboxEvent = ""                // перевод не выполнен
)

type AccountCreatedEvent struct {
	Id             uuid.UUID `json:"id"`
	InitialBalance float64   `json:"initial_balance"`
	CreatedAt      time.Time `json:"created_at"`
}

// Парсинг в payload как событие - создание нового аккаунта
func NewAccountCreatedEvent(event AccountCreatedEvent) (OutboxMessage, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return OutboxMessage{}, err
	}

	msg := NewOutboxMessage(payload, event.Id, OutboxEventAccountCreated)

	return msg, nil
}

type OutboxMessage struct {
	Id          uuid.UUID    `json:"id" db:"id"`
	AggregateId uuid.UUID    `json:"aggregate_id" db:"aggregate_id"`
	EventType   OutboxEvent  `json:"event_type" db:"event_type"`
	Payload     []byte       `json:"payload" db:"payload"`
	Status      OutboxStatus `json:"status" db:"status"`
	CreateAt    time.Time    `json:"created_at" db:"created_at"`
}

func NewOutboxMessage(payload []byte, aggregateId uuid.UUID, eventType OutboxEvent) OutboxMessage {
	return OutboxMessage{
		Id:          uuid.New(),
		AggregateId: aggregateId,
		EventType:   eventType,
		Payload:     payload,
		Status:      OutboxStatusPending,
		CreateAt:    time.Now().UTC().Truncate(time.Microsecond),
	}
}

type Payload struct {
}

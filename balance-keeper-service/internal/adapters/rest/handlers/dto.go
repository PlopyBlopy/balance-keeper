package handlers

import "github.com/google/uuid"

type addAccount struct {
	Id uuid.UUID `json:"id" binding:"require,uuid"`
}

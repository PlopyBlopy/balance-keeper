package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/PlopyBlopy/balance-keeper-service/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddAccount(u func(uuid.UUID, context.Context) (uuid.UUID, error)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var i addAccount

		err := ctx.ShouldBindJSON(&i)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := u(i.Id, ctx)
		if err != nil {
			if errors.Is(err, domain.ErrAlreadyExist) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func AddNewAccount() func(*gin.Context) {
	return func(ctx *gin.Context) {

	}
}

func GetAccount() func(*gin.Context) {
	return func(ctx *gin.Context) {

	}
}

func GetAccounts() func(*gin.Context) {
	return func(ctx *gin.Context) {

	}
}

func Deposit() func(*gin.Context) {
	return func(ctx *gin.Context) {

	}
}

func Withdraw() func(*gin.Context) {
	return func(ctx *gin.Context) {

	}
}

func Transfer() func(*gin.Context) {
	return func(ctx *gin.Context) {

	}
}

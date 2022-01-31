package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/ferueda/simplebank-go/db/sqlc"
	"github.com/ferueda/simplebank-go/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=CAD USD"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAcc, isValid := s.validateAccount(ctx, req.FromAccountId, req.Currency)
	if !isValid {
		return
	}

	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	if fromAcc.Owner != authPayload.Username {
		err := errors.New("wrong origin account")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	if !s.validateFunds(ctx, req.FromAccountId, req.Amount) {
		return
	}

	_, isValid = s.validateAccount(ctx, req.ToAccountId, req.Currency)
	if !isValid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Amount,
	}

	transfer, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, transfer)
}

type getTransferRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getTransfer(ctx *gin.Context) {
	var req getTransferRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := s.store.GetTransfer(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, transfer)
}

type listTransfersRequest struct {
	FromAccountId int64 `form:"from"`
	ToAccountId   int64 `form:"to"`
	Limit         int32 `form:"limit"`
	Offset        int32 `form:"offset"`
}

func (s *Server) listTransfers(ctx *gin.Context) {
	var req listTransfersRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.FromAccountId == 0 && req.ToAccountId == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("params for from or to account missing")))
		return
	}

	switch {
	case req.Limit <= 0:
		req.Limit = 20
	case req.Limit > 100:
		req.Limit = 100
	}

	switch {
	case req.Offset < 0:
		req.Offset = 0
	}

	arg := db.ListTransfersParams{
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountId,
		Limit:         req.Limit,
		Offset:        req.Offset,
	}

	transfers, err := s.store.ListTransfers(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"_metadata": map[string]interface{}{
			"count":  len(transfers),
			"offset": req.Offset,
		},
		"data": transfers,
	})
}

func (s *Server) validateAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: transaction must be in %s", accountId, account.Currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}

func (s *Server) validateFunds(ctx *gin.Context, accountId, amount int64) bool {
	account, err := s.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Balance < amount {
		err := fmt.Errorf("not enough funds in account [%d]", accountId)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}

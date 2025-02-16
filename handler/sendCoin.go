package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"merch-store/models"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

func (api *API) ApiSendCoin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username, err := validateAuthorizationHeader(r)
	if username == "" || err != nil {
		log.Printf("ERROR: Authorization failed. %v", err)
		err := sendErrorResponse(w, http.StatusUnauthorized, "Неавторизован.")
		if err != nil {
			log.Printf("ERROR: %v when Authorization failed", err)
		}
		return
	}
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read request body")
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос.")
		if err != nil {
			log.Printf("%v when failed to read request body", err)
		}

		return
	}

	var req models.SendCoinRequest
	if err = json.Unmarshal(buf.Bytes(), &req); err != nil {
		log.Printf("ERROR: invalid JSON format: %v", err)
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос. Некорректный формат JSON")
		if err != nil {
			log.Printf("%v when invalid JSON format", err)
		}

		return
	}

	if req.ToUser == "" {
		log.Printf("ERROR: ToUser cannot be empty: %v", err)
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос. ToUser должен быть заполнен")
		if err != nil {
			log.Printf("%v when toUser isn't filled", err)
		}

		return
	}
	if req.Amount <= 0 {
		log.Printf("ERROR: Amount should be more than zero: %v", err)
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос. Amount должен быть больше нуля")

		if err != nil {
			log.Printf("%v when amount is less than or equal to 0", err)
		}

		return
	}

	sender, err := api.store.GetUserCoinsByUserName(ctx, username)

	if err != nil || sender.UserID == 0 {
		log.Printf("ERROR: database error while fetching user coins: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")

		if err != nil {
			log.Printf("%v when fetching user coins", err)
		}

		return
	}

	if sender.Coins < int(req.Amount) {
		log.Printf("insufficient coins: available %d, requested %d", sender.Coins, req.Amount)
		err := sendErrorResponse(w, http.StatusPaymentRequired, "Недостаточно коинов")
		if err != nil {
			log.Printf("%v when insufficient coins", err)
		}
		return
	}

	recipient, err := api.store.GetUserCoinsByUserName(ctx, req.ToUser)
	//if err != nil || recipient.UserID == 0 {
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Printf("ERROR: database error while fetching recipient: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")

		if err != nil {
			log.Printf("%v fetching recipient", err)
		}

		return
	}

	if errors.Is(err, pgx.ErrNoRows) || recipient.UserID == 0 {
		log.Printf("ERROR: recipient not found: %v", err)
		err := sendErrorResponse(w, http.StatusNotFound, "Такой получатель не найден")
		if err != nil {
			log.Printf("ERROR: %v when recipient not found", err)
		}

		return
	}

	err = api.store.CreateTransactionSendCoin(ctx, models.Transaction{
		Type:        models.TransactionTypeSendCoins,
		Date:        time.Now(),
		SenderID:    sender.UserID,
		RecipientID: recipient.UserID,
		Coins:       int(req.Amount),
	})
	log.Println("transaction", err)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = sendResponse(w, models.APIResponse{
		StatusCode: http.StatusCreated,
		Headers: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
		Body: models.MessageResponse{
			Message: "Операция выполнена успешно",
		},
	})

	if err != nil {
		log.Printf("%v when trying to send a successful response.", err)
	}

}

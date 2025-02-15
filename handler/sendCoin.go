package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"merch-store/models"
	"net/http"
	"strings"
	"time"
)

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int32  `json:"amount"`
}

func (api *API) ApiSendCoin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Printf("ERROR: Authorization header required")
		err := sendErrorResponse(w, http.StatusUnauthorized, "Неавторизован.")
		if err != nil {
			log.Printf("%v when authorization header is not filled", err)
		}

		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		log.Printf("ERROR: Invalid authorization header format")
		err := sendErrorResponse(w, http.StatusUnauthorized, "Неавторизован.")
		if err != nil {
			log.Printf("%v when authorization header is not valid", err)
		}

		return
	}

	username, err := validateJWT(headerParts[1])
	if err != nil {
		log.Printf("ERROR: Invalid token")
		err := sendErrorResponse(w, http.StatusUnauthorized, "Неавторизован.")
		if err != nil {
			log.Printf("%v when Invalid token", err)
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

	var req SendCoinRequest
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
		log.Printf("ERROR: sender cannot be empty.: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")

		if err != nil {
			log.Printf("%v when sender isn't filled", err)
		}

		return
	}
	if sender.Coins < int(req.Amount) {
		log.Printf("insufficient coins: available %d, requested %d", sender.Coins, req.Amount)
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос. Недостаточно коинов")
		if err != nil {
			log.Printf("%v when insufficient coins", err)
		}
	}

	recipient, err := api.store.GetUserCoinsByUserName(ctx, req.ToUser)
	if err != nil || recipient.UserID == 0 {
		log.Printf("ERROR: recipient cannot be empty: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")

		if err != nil {
			log.Printf("%v when recipient isn't filled", err)
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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

}

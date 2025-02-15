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

type BuyItemRequest struct {
	Item string `json:"item"`
}

func (api *API) ApiBuyItem(w http.ResponseWriter, r *http.Request) {
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

	var req BuyItemRequest
	if err = json.Unmarshal(buf.Bytes(), &req); err != nil {
		log.Printf("ERROR: invalid JSON format: %v", err)
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос. Некорректный формат JSON")
		if err != nil {
			log.Printf("%v when invalid JSON format", err)
		}

		return
	}

	if req.Item == "" {
		log.Printf("ERROR: Item cannot be empty")
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос. Item должен быть заполнен")
		if err != nil {
			log.Printf("%v when Item is empty", err)
		}

		return
	}
	item, err := api.store.GetItemByName(ctx, req.Item)
	if err != nil {
		log.Printf("ERROR: invalid JSON format: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Неверный запрос. Некорректный формат JSON")
		if err != nil {
			log.Printf("%v when invalid JSON format", err)
		}

		return
	}
	buyer, err := api.store.GetUserCoinsByUserName(ctx, username)
	if err != nil || buyer.UserID == 0 {
		log.Printf("ERROR: buyer cannot be empty.: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")

		if err != nil {
			log.Printf("%v when buyer isn't filled", err)
		}

		return
	}
	if buyer.Coins < item.Price {
		log.Printf("insufficient coins: available %d, requested %d", buyer.Coins, item.Price)
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос. Недостаточно коинов")
		if err != nil {
			log.Printf("%v when insufficient coins", err)
		}
	}

	err = api.store.CreateTransactionBuyItem(ctx, models.Transaction{
		Type:        models.TransactionTypeBuyItem,
		Date:        time.Now(),
		SenderID:    buyer.UserID,
		RecipientID: buyer.UserID,
		Coins:       item.Price,
		ItemId:      item.ID,
	})
	log.Println("transaction", err)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

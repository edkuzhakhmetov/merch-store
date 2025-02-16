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

func (api *API) ApiBuyItem(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("ERROR: %v when failed to read request body", err)
		}

		return
	}

	var req models.BuyItemRequest
	if err = json.Unmarshal(buf.Bytes(), &req); err != nil {
		log.Printf("ERROR: invalid JSON format: %v", err)
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос. Некорректный формат JSON")
		if err != nil {
			log.Printf("ERROR: %v when invalid JSON format", err)
		}

		return
	}

	if req.Item == "" {
		log.Printf("ERROR: Item cannot be empty")
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос. Item должен быть заполнен")
		if err != nil {
			log.Printf("ERROR: %v when Item is empty", err)
		}

		return
	}
	item, err := api.store.GetItemByName(ctx, req.Item)
	if err != nil && (!errors.Is(err, pgx.ErrNoRows) || item.ID != 0) {
		log.Printf("ERROR: database error while fetching item: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		if err != nil {
			log.Printf("ERROR: %v when fetching item", err)
		}

		return
	}

	if errors.Is(err, pgx.ErrNoRows) || item.ID == 0 {
		log.Printf("ERROR: item not found: %v", err)
		err := sendErrorResponse(w, http.StatusNotFound, "Такой предмет не найден")
		if err != nil {
			log.Printf("ERROR: %v when item not found", err)
		}

		return
	}

	buyer, err := api.store.GetUserCoinsByUserName(ctx, username)
	if err != nil || buyer.UserID == 0 {
		log.Printf("ERROR: buyer cannot be empty.: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")

		if err != nil {
			log.Printf("ERROR: %v when buyer isn't filled", err)
		}

		return
	}
	if buyer.Coins < item.Price {
		log.Printf("insufficient coins: available %d, requested %d", buyer.Coins, item.Price)
		err := sendErrorResponse(w, http.StatusPaymentRequired, "Недостаточно коинов")
		if err != nil {
			log.Printf("ERROR: %v when insufficient coins", err)
		}
		return
	}

	err = api.store.CreateTransactionBuyItem(ctx, models.Transaction{
		Type:        models.TransactionTypeBuyItem,
		Date:        time.Now(),
		SenderID:    buyer.UserID,
		RecipientID: buyer.UserID,
		Coins:       item.Price,
		ItemId:      item.ID,
	})

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

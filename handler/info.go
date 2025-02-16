package handler

import (
	"log"
	"merch-store/models"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (api *API) GetInfo(w http.ResponseWriter, r *http.Request) {
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

	userCoins, err := api.store.GetUserCoinsByUserName(ctx, username)
	if err != nil || err == pgx.ErrNoRows {
		log.Printf("ERROR: user not found: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		if err != nil {
			log.Printf("%v when user not found", err)
		}

		return
	}

	inventory, err := api.store.GetUserInventoryByUserId(ctx, userCoins.UserID)
	if err != nil && err != pgx.ErrNoRows {
		log.Printf("ERROR: failed to get user inventory: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		if err != nil {
			log.Printf("%v when user not found", err)
		}

		return
	}
	coinHistory, err := api.store.GetCoinHistory(ctx, userCoins.UserID)
	if err != nil && err != pgx.ErrNoRows {
		log.Printf("ERROR: failed to get user inventory: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		if err != nil {
			log.Printf("%v when user not found", err)
		}

		return
	}

	err = sendResponse(w, models.APIResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
		Body: models.InfoResponse{
			Coins:       int32(userCoins.Coins),
			Inventory:   inventory,
			CoinHistory: *coinHistory,
		},
	})

	if err != nil {
		log.Printf("%v when trying to send a successful response.", err)
	}
}

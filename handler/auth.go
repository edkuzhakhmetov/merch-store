package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"merch-store/models"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

func (api *API) ApiAuth(w http.ResponseWriter, r *http.Request) {
	var user models.AuthRequest
	var buf bytes.Buffer

	ctx := r.Context()

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		log.Printf("ERROR: failed to read request body: %v", err)
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос.")
		if err != nil {
			log.Printf("ERROR: %v when failed to read request body", err)
		}
		return
	}

	err = json.Unmarshal(buf.Bytes(), &user)
	if err != nil {
		log.Printf("ERROR: invalid JSON format: %v", err)
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос. Некорректный формат JSON")
		if err != nil {
			log.Printf("%v when invalid JSON format", err)
		}
		return
	}

	if user.Username == "" || user.Password == "" {
		log.Println("ERROR: username and password are required")
		err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос. username и password должны быть заполнены")
		if err != nil {
			log.Printf("ERROR: %v when username or password are not filled", err)
		}
		return
	}

	internalUser, err := api.store.GetUserByUsername(ctx, user.Username)
	if err != nil && (!errors.Is(err, pgx.ErrNoRows) || internalUser.ID != 0) {
		log.Printf("ERROR: database error while fetching user or user not found: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		if err != nil {
			log.Printf("ERROR: %v while fetching user", err)
		}
		return
	}

	var isCreated bool

	if errors.Is(err, pgx.ErrNoRows) || internalUser.ID == 0 {
		salt := make([]byte, 16)
		_, err := rand.Read(salt)
		if err != nil {
			log.Printf("ERROR: failed to generate salt: %v", err)
			err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
			if err != nil {
				log.Printf("ERROR: %v when failed to generate salt", err)
			}
			return
		}

		hashedPassword := sha256.Sum256([]byte(user.Password + string(salt)))
		internalUser, err = api.store.CreateUser(ctx, models.User{
			Username:       user.Username,
			HashedPassword: hex.EncodeToString(hashedPassword[:]),
			Salt:           hex.EncodeToString(salt),
			CreatedAt:      time.Now(),
		})
		if err != nil {
			log.Printf("ERROR: failed to create user: %v", err)
			err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
			if err != nil {
				log.Printf("ERROR: %v when failed to create user", err)
			}
			return
		}
		isCreated = true
	}

	salt, err := hex.DecodeString(internalUser.Salt)
	if err != nil {
		log.Printf("ERROR: failed to decode salt for user verification: %v", err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		if err != nil {
			log.Printf("ERROR: %v when failed to decode salt for user verification", err)
		}
		return
	}

	hashedInputPassword := sha256.Sum256([]byte(user.Password + string(salt)))
	hashedInputPasswordHex := hex.EncodeToString(hashedInputPassword[:])

	if !hmac.Equal([]byte(hashedInputPasswordHex), []byte(internalUser.HashedPassword)) {
		log.Printf("ERROR: Invalid username or password for user: %s", user.Username)
		err := sendErrorResponse(w, http.StatusUnauthorized, "Неавторизован.")
		if err != nil {
			log.Printf("ERROR: %v when Invalid username or password", err)
		}
		return
	}

	token, err := generateJWT(user.Username)
	if err != nil {
		log.Printf("ERROR: Failed to generate JWT for user: %s, error: %v", user.Username, err)
		err := sendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		if err != nil {
			log.Printf("ERROR: %v when failed to generate JWT for user", err)
		}
		return
	}

	status := http.StatusOK
	if isCreated {
		status = http.StatusCreated
	}

	err = sendResponse(w, models.APIResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
		Body: models.AuthResponse{
			Token: token,
		},
	})

	if err != nil {
		log.Printf("%v when trying to send a successful response.", err)
	}
}

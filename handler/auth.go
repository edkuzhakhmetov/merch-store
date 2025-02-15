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

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
}

func (api *API) ApiAuth(w http.ResponseWriter, r *http.Request) {
	var user AuthRequest
	var buf bytes.Buffer

	ctx := r.Context()

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		log.Printf("ERROR: failed to read request body: %v", err)
		err := sendResponse(w, APIResponse{
			StatusCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json; charset=UTF-8",
			},
			Body: ErrorResponse{
				Errors: "Неверный запрос.",
			},
		})
		if err != nil {
			log.Printf("%v when failed to read request body", err)
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
		err = sendResponse(w, APIResponse{
			StatusCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json; charset=UTF-8",
			},
			Body: ErrorResponse{
				Errors: "Неверный запрос. username и password должны быть заполнены",
			},
		})

		if err != nil {
			log.Printf("%v when username or password are not filled", err)
		}
		return
	}

	internalUser, err := api.store.GetUserByUsername(ctx, user.Username)
	if err != nil && (!errors.Is(err, pgx.ErrNoRows) || internalUser.ID != 0) {
		log.Printf("ERROR: database error while fetching user: %v", err)
		err = sendResponse(w, APIResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type": "application/json; charset=UTF-8",
			},
			Body: ErrorResponse{
				Errors: "Внутренняя ошибка сервера",
			},
		})

		if err != nil {
			log.Printf("%v while fetching user", err)
		}
		return
	}

	if errors.Is(err, pgx.ErrNoRows) || internalUser.ID == 0 {
		salt := make([]byte, 16)
		_, err := rand.Read(salt)
		if err != nil {
			log.Printf("ERROR: failed to generate salt: %v", err)
			err = sendResponse(w, APIResponse{
				StatusCode: http.StatusInternalServerError,
				Headers: map[string]string{
					"Content-Type": "application/json; charset=UTF-8",
				},
				Body: ErrorResponse{
					Errors: "Внутренняя ошибка сервера",
				},
			})

			if err != nil {
				log.Printf("%v when failed to generate salt", err)
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
			err = sendResponse(w, APIResponse{
				StatusCode: http.StatusInternalServerError,
				Headers: map[string]string{
					"Content-Type": "application/json; charset=UTF-8",
				},
				Body: ErrorResponse{
					Errors: "Внутренняя ошибка сервера",
				},
			})

			if err != nil {
				log.Printf("%v when failed to create user", err)
			}
			return
		}

	}

	salt, err := hex.DecodeString(internalUser.Salt)
	if err != nil {
		log.Printf("ERROR: failed to decode salt for user verification: %v", err)
		err = sendResponse(w, APIResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type": "application/json; charset=UTF-8",
			},
			Body: ErrorResponse{
				Errors: "Внутренняя ошибка сервера",
			},
		})

		if err != nil {
			log.Printf("%v when failed to decode salt for user verification", err)
		}
		return
	}

	hashedInputPassword := sha256.Sum256([]byte(user.Password + string(salt)))
	hashedInputPasswordHex := hex.EncodeToString(hashedInputPassword[:])

	if !hmac.Equal([]byte(hashedInputPasswordHex), []byte(internalUser.HashedPassword)) {
		log.Printf("ERROR: Invalid username or password for user: %s", user.Username)
		err = sendResponse(w, APIResponse{
			StatusCode: http.StatusUnauthorized,
			Headers: map[string]string{
				"Content-Type": "application/json; charset=UTF-8",
			},
			Body: ErrorResponse{
				Errors: "Неавторизован.",
			},
		})

		if err != nil {
			log.Printf("%v when Invalid username or password", err)
		}
		return
	}

	token, err := generateJWT(user.Username)
	if err != nil {
		log.Printf("ERROR: Failed to generate JWT for user: %s, error: %v", user.Username, err)
		err = sendResponse(w, APIResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type": "application/json; charset=UTF-8",
			},
			Body: ErrorResponse{
				Errors: "Неавторизован.",
			},
		})

		if err != nil {
			log.Printf("ERROR: Failed to send successful response for user: %s, error: %v", user.Username, err)
		}
		return
	}

	err = sendResponse(w, APIResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
		Body: AuthResponse{
			Token: token,
		},
	})

	if err != nil {
		log.Printf("%v when trying to send a successful response.", err)
	}
}

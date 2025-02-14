package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"merch-store/config"
	"merch-store/models"
	"merch-store/storage"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
}

type ErrorResponse struct {
	Errors string `json:"errors,omitempty"`
}

var secret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func ApiAuth(w http.ResponseWriter, r *http.Request) {
	var user AuthRequest
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		response, _ := json.Marshal(ErrorResponse{
			Errors: "Неверный запрос",
		})

		//http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	response, err := json.Marshal(AuthResponse{
		Token: "token",
	})
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func main() {
	log.SetFlags(5)
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load config. %v", err)
	}
	log.Println("config ok")

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, time.Millisecond*500)
	defer cancel()
	//ctx := context.TODO()

	db, err := storage.Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close(ctx)
	log.Println("Connected to database")

	pg := storage.NewStore(db)

	user, err := pg.GetUserByUsername(ctx, "ed")
	log.Println(user, err)

	newOldUser, err := pg.CreateUser(ctx, user)
	log.Println(newOldUser, err)

	newUser, err := pg.CreateUser(ctx, models.User{
		Username:       "admin1",
		HashedPassword: "pas",
		Salt:           "salt",
	})
	fmt.Println(newUser, err)
	coins, err := pg.GetUserCoins(ctx, newUser.ID)
	log.Println(coins, err)

	r := chi.NewRouter()
	r.Post("/api/auth", ApiAuth)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed ListenAndServe. %v", err.Error())
		return
	}
}

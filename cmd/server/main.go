package main

import (
	"context"
	"log"
	"merch-store/config"
	"merch-store/handler"
	"merch-store/storage"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func main() {
	log.SetFlags(5)
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("FATAL ERROR: Failed to load config. %v", err)
	}

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, time.Millisecond*500)
	defer cancel()
	//ctx := context.TODO()

	db, err := storage.Connect(ctx)
	if err != nil {
		log.Fatalf("FATAL ERROR: Failed to connect to database: %v", err)
	}
	defer db.Close(ctx)

	pg := storage.NewStore(db)

	us, err := pg.GetUserItemsByUserName(ctx, "ed_6")
	log.Println(us, err)
	r := chi.NewRouter()

	api := handler.NewApi(pg)

	r.Post("/api/auth", api.ApiAuth)
	r.Post("/api/sendCoin", api.ApiSendCoin)
	r.Post("/api/buyItem", api.ApiBuyItem)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("FATAL ERROR: Failed ListenAndServe. %v", err.Error())

		return
	}

}

package handler

import (
	"merch-store/storage"
)

type API struct {
	store *storage.StorePostgres
}

func NewApi(store *storage.StorePostgres) *API {
	return &API{store: store}
}

package storage

import (
	"context"
	"merch-store/models"
)

func (s *StorePostgres) GetItemByName(ctx context.Context, itemName string) (models.Item, error) {
	var item models.Item
	err := s.db.QueryRow(ctx, "SELECT id, name, price FROM merch.items WHERE name = $1", itemName).Scan(&item.ID, &item.Name, &item.Price)
	return item, err
}

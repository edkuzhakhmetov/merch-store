package storage

import (
	"context"
	"fmt"
	"merch-store/models"
)

func (s *StorePostgres) GetUserItemsByUserName(ctx context.Context, userName string) (models.UserItem2, error) {

	var ui models.UserItem2
	err := s.db.QueryRow(ctx,
		`SELECT u.id, i.Name, i.id, ui.quantity
		 FROM merch.users u
		 JOIN merch.user_items ui ON u.id=ui.user_id
		 JOIN merch.items i ON ui.item_id=i.id
		 WHERE u.username= $1`,
		userName).Scan(&ui.User.ID, &ui.Item.Name, &ui.Item.ID, &ui.Quantity)
	if err != nil {
		return models.UserItem2{}, fmt.Errorf("failed to get user items: %w", err)
	}
	return ui, nil
}

func (s *StorePostgres) GetCoinHistoryReceived(ctx context.Context, userId int) (models.InfoResponseCoinHistoryReceived, error) {

	var hist models.InfoResponseCoinHistoryReceived
	err := s.db.QueryRow(ctx,
		`WITH t as (SELECT t.sender_id, SUM(t.coins) coins
		FROM merch.transactions t
		WHERE t.type = $1
		AND t.recipient_id = $2
		GROUP BY t.sender_id
		)
		SELECT u.username, t.coins
		FROM t 
		JOIN merch.users u on u.id=t.sender_id`,
		models.TransactionTypeSendCoins, userId).Scan(&hist.FromUser, &hist.Amount)
	if err != nil {
		return models.InfoResponseCoinHistoryReceived{}, fmt.Errorf("failed to get user items: %w", err)
	}
	return hist, nil
}

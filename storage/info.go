package storage

import (
	"context"
	"fmt"
	"merch-store/models"

	"github.com/jackc/pgx/v5"
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

func (s *StorePostgres) GetUserInventoryByUserId(ctx context.Context, userId int) ([]models.InfoResponseInventory, error) {

	var res []models.InfoResponseInventory
	rows, err := s.db.Query(ctx,
		`SELECT u.id, i.Name, ui.quantity
		 FROM merch.users u
		 JOIN merch.user_items ui ON u.id=ui.user_id
		 JOIN merch.items i ON ui.item_id=i.id
		 WHERE u.id=$1`,
		userId)

	defer rows.Close()

	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user item rows: %w", err)
	}

	for rows.Next() {
		row := models.InfoResponseInventory{}
		err = rows.Scan(&row.User.ID, &row.Type, &row.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to get user item row: %w", err)
		}
		res = append(res, row)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("an error occurred while fetching user item rows %w", err)
	}

	return res, nil

}

func (s *StorePostgres) GetCoinHistory(ctx context.Context, userId int) (*models.InfoResponseCoinHistory, error) {
	received, err := s.GetCoinHistoryReceived(ctx, userId)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to get received coins history: %w", err)
	}

	sent, err := s.GetCoinHistorySent(ctx, userId)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to get sent coins history: %w", err)
	}

	return &models.InfoResponseCoinHistory{
		Received: received,
		Sent:     sent,
	}, nil
}

func (s *StorePostgres) GetCoinHistoryReceived(ctx context.Context, userId int) ([]models.InfoResponseCoinHistoryReceived, error) {

	var res []models.InfoResponseCoinHistoryReceived
	rows, err := s.db.Query(ctx,
		`WITH t AS
			(SELECT t.sender_id, SUM(t.coins) coins
				FROM merch.transactions t
				WHERE t.type = $1
				AND t.recipient_id = $2
				GROUP BY t.sender_id
			)
		SELECT u.username, t.coins
		FROM t
		JOIN merch.users u on u.id=t.sender_id`,
		models.TransactionTypeSendCoins, userId)
	defer rows.Close()
	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user item rows: %w", err)
	}

	for rows.Next() {
		row := models.InfoResponseCoinHistoryReceived{}
		err = rows.Scan(&row.FromUser, &row.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to get user item row: %w", err)
		}
		res = append(res, row)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("an error occurred while fetching user item rows %w", err)
	}

	return res, nil
}

/*

 */

func (s *StorePostgres) GetCoinHistorySent(ctx context.Context, userId int) ([]models.InfoResponseCoinHistorySent, error) {
	var res []models.InfoResponseCoinHistorySent
	rows, err := s.db.Query(ctx,
		`WITH t AS
			(SELECT t.recipient_id, SUM(t.coins) coins
				FROM merch.transactions t
				WHERE t.type=$1
				AND t.sender_id=$2 
				GROUP BY t.recipient_id
			)
		SELECT u.username, t.coins
		FROM t
		JOIN merch.users u on u.id=t.recipient_id`,
		models.TransactionTypeSendCoins, userId)

	defer rows.Close()
	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user item rows: %w", err)
	}

	for rows.Next() {
		row := models.InfoResponseCoinHistorySent{}
		err = rows.Scan(&row.ToUser, &row.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to get user item row: %w", err)
		}
		res = append(res, row)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("an error occurred while fetching user item rows %w", err)
	}

	return res, nil
}

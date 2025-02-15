package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"merch-store/models"

	"github.com/jackc/pgx/v5"
)

func (s *StorePostgres) CreateTransactionSendCoin(ctx context.Context, transaction models.Transaction) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var senderCoins models.UserCoin

	err = tx.QueryRow(ctx, `
	UPDATE merch.user_coins 
	SET coins = coins - $1 
	WHERE user_id = $2 AND coins >= $1 
	RETURNING user_id, coins`,
		transaction.Coins, transaction.SenderID).Scan(&senderCoins.UserID, &senderCoins.Coins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("insufficient coins or sender user not found (user_id: %d)", transaction.SenderID)
		}
		return fmt.Errorf("failed to update sender balance: %w", err)
	}

	var recipientCoins models.UserCoin
	err = tx.QueryRow(ctx, `
	UPDATE merch.user_coins 
	SET coins = coins + $1 
	WHERE user_id = $2 
	RETURNING user_id, coins`,
		transaction.Coins, transaction.RecipientID).Scan(&recipientCoins.UserID, &recipientCoins.Coins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("recipient user not found (user_id: %d)", transaction.RecipientID)
		}
		return fmt.Errorf("failed to update recipient balance: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO merch.transactions
		 (type, created_at, sender_id, recipient_id, coins)
		 VALUES ($1, $2, $3, $4, $5)`,
		transaction.Type, transaction.Date, transaction.SenderID, transaction.RecipientID, transaction.Coins)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	log.Printf("transaction successfully created from user %d to user %d for %d coins", transaction.SenderID, transaction.RecipientID, transaction.Coins)

	return nil
}

func (s *StorePostgres) CreateTransactionBuyItem(ctx context.Context, transaction models.Transaction) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var buyerCoins models.UserCoin

	err = tx.QueryRow(ctx, `
	UPDATE merch.user_coins 
	SET coins = coins - $1 
	WHERE user_id = $2 AND coins >= $1 
	RETURNING user_id, coins`,
		transaction.Coins, transaction.SenderID).Scan(&buyerCoins.UserID, &buyerCoins.Coins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("insufficient coins or buyer user not found (user_id: %d)", transaction.SenderID)
		}
		return fmt.Errorf("failed to update buyer balance: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO merch.user_items (user_id, item_id, quantity) 
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, item_id) 
		DO UPDATE SET quantity = user_items.quantity + $3`,
		transaction.RecipientID, transaction.ItemId, 1)
	if err != nil {
		return fmt.Errorf("failed to update user items: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO merch.transactions (type, created_at, recipient_id, coins, item_id) 
		VALUES ($1, $2, $3, $4, $5)`,
		transaction.Type, transaction.Date, transaction.RecipientID, -transaction.Coins, transaction.ItemId)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx.Commit(ctx)
}

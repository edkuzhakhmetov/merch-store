package storage

import (
	"context"
	"errors"
	"fmt"
	"merch-store/models"
	"time"

	"github.com/jackc/pgx/v5"
)

const initialCoins = 1000

func (s *StorePostgres) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	err := s.db.QueryRow(ctx, "SELECT id, username, hashed_password, salt, created_at FROM merch.users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Salt, &user.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.User{}, pgx.ErrNoRows
	}

	if err != nil {
		return models.User{}, fmt.Errorf("error retrieving user %s: %w", username, err)
	}

	return user, nil
}

func (s *StorePostgres) CreateUser(ctx context.Context, user models.User) (models.User, error) {

	foundUser, err := s.GetUserByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, fmt.Errorf("failed to check existing user: %w", err)
	}
	if foundUser.ID != 0 {
		return models.User{}, fmt.Errorf("user with username '%s' already exists", user.Username)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var newUser models.User
	err = tx.QueryRow(ctx,
		`INSERT INTO merch.users (username, hashed_password, salt, created_at) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id, username, hashed_password, salt, created_at`,
		user.Username, user.HashedPassword, user.Salt, time.Now(),
	).Scan(&newUser.ID, &newUser.Username, &newUser.HashedPassword, &newUser.Salt, &newUser.CreatedAt)

	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	if _, err = tx.Exec(ctx,
		`INSERT INTO merch.user_coins (user_id, coins) VALUES ($1, $2)`,
		newUser.ID, initialCoins,
	); err != nil {
		return models.User{}, fmt.Errorf("failed to set user coins: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return models.User{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newUser, nil
}

func (s *StorePostgres) GetUserCoinsByUserID(ctx context.Context, userID int) (int, error) {
	var coins int

	err := s.db.QueryRow(ctx, "SELECT coins FROM merch.user_coins WHERE user_id = $1", userID).Scan(&coins)

	return coins, err
}

func (s *StorePostgres) GetUserCoinsByUserName(ctx context.Context, userName string) (models.UserCoin, error) {
	var userCoin models.UserCoin

	err := s.db.QueryRow(ctx,
		`SELECT u.id, COALESCE(c.coins, 0) AS coins
		 FROM merch.users u
		 LEFT JOIN merch.user_coins c ON u.id=c.user_id
		 WHERE u.username= $1`,
		userName).Scan(&userCoin.UserID, &userCoin.Coins)
	if err != nil {
		return models.UserCoin{}, fmt.Errorf("failed to get user coins: %w", err)
	}
	return userCoin, err
}

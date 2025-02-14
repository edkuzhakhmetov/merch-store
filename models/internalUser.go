package models

import (
	"time"
)

type User struct {
	ID             int
	Username       string
	HashedPassword string
	Salt           string
	CreatedAt      time.Time
}

type UserCoin struct {
	UserID int
	Coins  int
}

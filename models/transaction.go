package models

import "time"

const (
	TransactionTypeSendCoins = 1
	TransactionTypeBuyItem   = 2
)

type Transaction struct {
	ID          int
	Type        int
	Date        time.Time
	SenderID    int
	RecipientID int
	Coins       int
	ItemId      int
}

type InfoResponseCoinHistorySent struct {
	ToUser string `json:"toUser,omitempty"`
	Amount int32  `json:"amount,omitempty"`
}

type InfoResponseCoinHistoryReceived struct {
	FromUser string `json:"fromUser,omitempty"`
	Amount   int32  `json:"amount,omitempty"`
}

type InfoResponseCoinHistory struct {
	Received []InfoResponseCoinHistoryReceived `json:"received,omitempty"`
	Sent     []InfoResponseCoinHistorySent     `json:"sent,omitempty"`
}

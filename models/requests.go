package models

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int32  `json:"amount"`
}

type BuyItemRequest struct {
	Item string `json:"item"`
}

type ErrorResponse struct {
	Errors string `json:"errors,omitempty"`
}

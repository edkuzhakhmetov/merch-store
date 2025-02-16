package models

type AuthResponse struct {
	Token string `json:"token,omitempty"`
}

type APIResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       interface{}
}

type MessageResponse struct {
	Message string `json:"message,omitempty"`
}

type InfoResponseInventory struct {
	User     User   `json:"-"`
	Type     string `json:"type,omitempty"`
	Quantity int32  `json:"quantity,omitempty"`
}

type InfoResponse struct {
	Coins       int32                   `json:"coins,omitempty"`
	Inventory   []InfoResponseInventory `json:"inventory,omitempty"`
	CoinHistory InfoResponseCoinHistory `json:"coinHistory,omitempty"`
}

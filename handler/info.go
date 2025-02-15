package handler

import "merch-store/models"

type InfoResponse struct {
	// Количество доступных монет.
	Coins       int32                           `json:"coins,omitempty"`
	Inventory   []InfoResponseInventory         `json:"inventory,omitempty"`
	CoinHistory *models.InfoResponseCoinHistory `json:"coinHistory,omitempty"`
}

type InfoResponseInventory struct {
	Type_    string `json:"type,omitempty"`
	Quantity int32  `json:"quantity,omitempty"`
}

/* func (api *API) GetInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Printf("ERROR: Authorization header required")
		err := sendErrorResponse(w, http.StatusUnauthorized, "Неавторизован.")
		if err != nil {
			log.Printf("%v when authorization header is not filled", err)
		}

		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		log.Printf("ERROR: Invalid authorization header format")
		err := sendErrorResponse(w, http.StatusUnauthorized, "Неавторизован.")
		if err != nil {
			log.Printf("%v when authorization header is not valid", err)
		}

		return
	}

	username, err := validateJWT(headerParts[1])
	if err != nil {
		log.Printf("ERROR: Invalid token")
		err := sendErrorResponse(w, http.StatusUnauthorized, "Неавторизован.")
		if err != nil {
			log.Printf("%v when Invalid token", err)
		}

		return
	}

	// var buf bytes.Buffer

	// _, err = buf.ReadFrom(r.Body)
	// if err != nil {
	// 	log.Printf("ERROR: Failed to read request body")
	// 	err := sendErrorResponse(w, http.StatusBadRequest, "Неверный запрос.")
	// 	if err != nil {
	// 		log.Printf("%v when failed to read request body", err)
	// 	}

	// 	return
	// }

	/*
	   type InfoResponseInventory struct {
	   	// Тип предмета.
	   	Type_ string `json:"type,omitempty"`
	   	// Количество предметов.
	   	Quantity int32 `json:"quantity,omitempty"`

	   type UserItem struct {
	   	UserID   int
	   	ItemId   int
	   	Quantity int
	   }



}
*/

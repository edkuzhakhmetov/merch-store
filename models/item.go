package models

type Item struct {
	ID    int
	Name  string
	Price int
}

// type UserItem struct {
// 	UserID   int
// 	ItemId   int
// 	Quantity int
// }

type UserItem2 struct {
	User     User
	Item     Item
	Quantity int
}

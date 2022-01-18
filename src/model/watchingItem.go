package model

type WatchingItem struct {
	ItemID        string `json:"id" gorm:"primary_key"`
	UserReference string
	Url           string  `json:"url"`
	ItemName      string  `json:"item_name"`
	LastChecked   string  `json:"last_check"`
	NextCheck     string  `json:"next_check"`
	CurrentPrice  float64 `json:"curr_price"`
}

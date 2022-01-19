package model

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model

	UserID       uint
	Url          string  `json:"url"`
	ItemName     string  `json:"item_name"`
	LastChecked  string  `json:"last_check"`
	NextCheck    string  `json:"next_check"`
	CurrentPrice float64 `json:"curr_price"`
}

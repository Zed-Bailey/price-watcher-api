package model

type User struct {
	UserID           string `gorm:"primary_key"`
	Email            string
	Password         string
	WatchingProducts []WatchingItem `json:"items" gorm:"foreignKey:UserReference"`
}

package model

type DateSlot struct {
	Id int `json:"id" db:"id" `
	Date string `json:"date" db:"date"`
}
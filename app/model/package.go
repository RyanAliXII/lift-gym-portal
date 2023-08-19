package model

type Package struct {
	Id          int     `db:"id" json:"id"`
Description string  `db:"description" json:"description"`
	Price       float64 `db:"price" json:"price"`
}
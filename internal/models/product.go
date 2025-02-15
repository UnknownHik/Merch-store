package models

type Product struct {
	Item  string `db:"item"`
	Price int    `db:"price"`
}

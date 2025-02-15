package models

type User struct {
	UserName string `db:"username"`
	Password string `db:"password"`
	Balance  int    `db:"balance"`
}

package dto

// UserAuth представляет данные для аутентификации пользователя
type UserAuth struct {
	UserName string `json:"username" binding:"required,username"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

// SendCoin представляет данные для отправки монет
type SendCoin struct {
	ToUser string `json:"toUser" binding:"required,username"`
	Amount int    `json:"amount" binding:"required,min=1"`
}

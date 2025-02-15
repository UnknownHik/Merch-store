package dto

// InfoResponse представляет сводные данные о балансе и действиях пользователя
type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []Item      `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

// Item представляет данные о приобретенном товаре
type Item struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

// CoinHistory представляет данные о транзакциях пользователя
type CoinHistory struct {
	Received []ReceivedCoin `json:"received,omitempty"`
	Sent     []SentCoin     `json:"sent,omitempty"`
}

// ReceivedCoin представляет данные от кого были получены монеты
type ReceivedCoin struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

// SentCoin представляет данные кому были отправлены монеты
type SentCoin struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

package models

type Wallet struct {
	WalletID string `json:"wallet_id" db:"id"`
	Balance  int32  `json:"balance" db:"balance"`
}

package models

type Transaction struct {
	ID               int64  `json:"id"`
	SenderWalletID   int64  `json:"sender_wallet_id"`
	ReceiverWalletID int64  `json:"receiver_wallet_id"`
	Amount           int64  `json:"amount"`
	CreatedAt        string `json:"created_at"`
}

type LedgerEntry struct {
	ID            int64  `json:"id"`
	TransactionID int64  `json:"transaction_id"`
	WalletID      int64  `json:"wallet_id"`
	EntryType     string `json:"entry_type"` // debit or credit
	Amount        int64  `json:"amount"`
	CreatedAt     string `json:"created_at"`
}

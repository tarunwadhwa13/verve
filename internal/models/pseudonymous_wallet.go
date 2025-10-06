package models

type PseudonymousWallet struct {
	ID             int64  `json:"id"`
	PublicKey      string `json:"public_key"`
	IsPseudonymous bool   `json:"is_pseudonymous"`
}

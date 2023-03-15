package iron

type ViewAccount struct {
	Version     int    `json:"version"`
	Name        string `json:"name"`
	PublicKey   string `json:"publicKey"`
	ViewKey     string `json:"viewKey"`
	OutgoingKey string `json:"outgoingKey"`
	IncomingKey string `json:"incomingKey"`
	CreatedAt   string `json:"createdAt"`
}

type SignTxMsg struct {
	FromAccount string `json:"fromAccount"`
	Transaction string `json:"transaction"`
}

type BroadcastTxMsg struct {
	FromAccount       string `json:"fromAccount"`
	SignedTransaction string `json:"signedTransaction"`
}

type SyncTxMsg struct {
	FromAccount string `json:"fromAccount"`
	TxHash      string `json:"tx_id"` //for updateTransaction
}

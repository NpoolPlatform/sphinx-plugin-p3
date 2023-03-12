package iron

import (
	ct "github.com/NpoolPlatform/sphinx-plugin-p3/pkg/types"
)

type IronfishAccount struct {
	PublicKey   string `json:"publicKey"`
	ViewKey     string `json:"viewKey"`
	OutgoingKey string `json:"outgoingKey"`
	IncomingKey string `json:"incomingKey"`
}

type SignMsgTx struct {
	BaseInfo        ct.BaseInfo `json:"base_info"`
	RecentBlockHash string      `json:"recent_block_hash"`
}

type BroadcastRequest struct {
	Signature []byte `json:"signature"`
}

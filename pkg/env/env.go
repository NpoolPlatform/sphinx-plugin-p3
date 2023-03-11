package env

import (
	"errors"
	"os"
)

const (
	// main or test
	ENVCOINNET = "ENV_COIN_NET"

	// ENVSYNCINTERVAL sync transaction status on chain interval
	ENVSYNCINTERVAL = "ENV_SYNC_INTERVAL"

	// for all chain
	ENVCOINTYPE      = "ENV_COIN_TYPE"
	ENVCOINLOCALAPI  = "ENV_COIN_LOCAL_API"
	ENVCOINPUBLICAPI = "ENV_COIN_PUBLIC_API"

	// for tokens
	ENVCONTRACT = "ENV_CONTRACT"

	ENVBUIILDCHIANSERVER = "ENV_BUILD_CHAIN_SERVER"
	ENVPROXY             = "ENV_PROXY"
)

var (
	// env error----------------------------
	ErrEVNCoinType     = errors.New("env ENV_COIN_TYPE not found")
	ErrEVNCoinNet      = errors.New("env ENV_COIN_NET not found")
	ErrEVNCoinNetValue = errors.New("env ENV_COIN_NET value only support main|test")

	ErrENVCoinLocalAPINotFound  = errors.New("env ENV_COIN_LOCAL_API not found")
	ErrENVCoinPublicAPINotFound = errors.New("env ENV_COIN_PUBLIC_API not found")

	// eth/usdt
	ErrENVContractNotFound         = errors.New("env ENV_CONTRACT not found")
	ErrENVBuildChainServerNotFound = errors.New("env ENV_BUILD_CHAIN_SERVER not found")
	ErrENVBuildChainServerInvalid  = errors.New("env ENV_BUILD_CHAIN_SERVER invalid")
	ErrENVProxyInvalid             = errors.New("env ENV_PROXY invalid")

	// tron
	ErrENVCOINJSONRPCAPINotFound = errors.New("env ENV_COIN_JSONRPC_API not found")
	ErrENVCOINGRPCAPINotFound    = errors.New("env ENV_COIN_GRPC_API not found")

	// not env error----------------------------
	ErrSignTypeInvalid     = errors.New("sign type invalid")
	ErrFindMsgNotFound     = errors.New("failed to find message")
	ErrCIDInvalid          = errors.New("cid invalid")
	ErrAddressInvalid      = errors.New("address invalid")
	ErrAmountInvalid       = errors.New("amount invalid")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrWaitMessageOnChain  = errors.New("wait message on chain")
	ErrContractInvalid     = errors.New("invalid contract address")
	ErrTransactionFail     = errors.New("transaction fail")
)

func LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func CoinInfo() (networkType, coinType string, err error) {
	var ok bool
	networkType, ok = LookupEnv(ENVCOINNET)
	if !ok {
		err = ErrEVNCoinNet
		return
	}

	coinType, ok = LookupEnv(ENVCOINTYPE)
	if !ok {
		err = ErrEVNCoinType
		return
	}
	return
}

package iron

import (
	"errors"
	"math/big"
	"strings"

	"github.com/NpoolPlatform/message/npool/sphinxplugin"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/coins"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/coins/register"
)

const (
	IronPrePoint = 100000000
)

var (
	// EmptyWalletL ..
	EmptyWalletL = big.Int{}
	// EmptyWalletS ..
	EmptyWalletS = big.Float{}
)

var (
	// ErrNodeNotSynced ..
	ErrNodeNotSynced = errors.New("node not synced or stoped")
	// ErrConnotGetBalance ..
	ErrConnotGetBalance = errors.New("cannot get balance from iron node")
	// ErrImportWalletWrong ..
	ErrImportWalletWrong = errors.New("import wallet failed")
)

var (
	SolTransactionFailed = `sol transaction failed`
	lamportsLow          = `Transfer: insufficient lamports`
	stopErrMsg           = []string{lamportsLow, SolTransactionFailed}
	ironfishToken        = &coins.TokenInfo{OfficialName: "IronFish", Decimal: 8, Unit: "IRON", Name: "ironfish", OfficialContract: "ironfish", TokenType: coins.Ironfish}
)

func init() {
	ironfishToken.Waight = 100
	ironfishToken.Net = coins.CoinNetMain
	ironfishToken.Contract = ironfishToken.OfficialContract
	ironfishToken.CoinType = sphinxplugin.CoinType_CoinTypeironfish
	register.RegisteTokenInfo(ironfishToken)
}

func ToIron(point uint64) *big.Float {
	// Convert lamports to sol:
	return big.NewFloat(0).
		Quo(
			big.NewFloat(0).SetUint64(point),
			big.NewFloat(0).SetUint64(IronPrePoint),
		)
}

func ToPoint(value float64) (uint64, big.Accuracy) {
	return big.NewFloat(0).Mul(
		big.NewFloat(0).SetFloat64(value),
		big.NewFloat(0).SetUint64(IronPrePoint),
	).Uint64()
}

func TxFailErr(err error) bool {
	if err == nil {
		return false
	}

	for _, v := range stopErrMsg {
		if strings.Contains(err.Error(), v) {
			return true
		}
	}
	return false
}

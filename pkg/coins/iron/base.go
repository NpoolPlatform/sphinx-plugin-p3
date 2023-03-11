package iron

import (
	"errors"
	"math/big"
	"strings"

	"github.com/NpoolPlatform/message/npool/sphinxplugin"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/coins"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/coins/register"
	solana "github.com/gagliardetto/solana-go"
)

var (
	// EmptyWalletL ..
	EmptyWalletL = big.Int{}
	// EmptyWalletS ..
	EmptyWalletS = big.Float{}
)

var (
	// ErrSolBlockNotFound ..
	ErrSolBlockNotFound = errors.New("not found confirmed block in solana chain")
	// ErrSolSignatureWrong ..
	ErrSolSignatureWrong = errors.New("solana signature is wrong or failed")
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

func ToSol(larm uint64) *big.Float {
	// Convert lamports to sol:
	return big.NewFloat(0).
		Quo(
			big.NewFloat(0).SetUint64(larm),
			big.NewFloat(0).SetUint64(solana.LAMPORTS_PER_SOL),
		)
}

func ToLarm(value float64) (uint64, big.Accuracy) {
	return big.NewFloat(0).Mul(
		big.NewFloat(0).SetFloat64(value),
		big.NewFloat(0).SetUint64(solana.LAMPORTS_PER_SOL),
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

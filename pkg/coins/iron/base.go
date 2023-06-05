package iron

import (
	"errors"
	"math/big"
	"strings"

	v1 "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/NpoolPlatform/message/npool/sphinxplugin"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/coins"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/coins/register"
	"github.com/shopspring/decimal"
)

const (
	IronExp                    = -8
	DefaultConfirmations       = 2
	ToleranceHeight            = 1
	TransactionExpirationDelta = 15
	// $IRON 0.0001
	MaxFeeLimit = 10000

	ChainType           = sphinxplugin.ChainType_Ironfish
	ChainNativeUnit     = "IRON"
	ChainAtomicUnit     = "ORE"
	ChainUnitExp        = 8
	ChainID             = "1"
	ChainNickname       = "Ironfish"
	ChainNativeCoinName = "ironfish"
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
	// ErrAccountNotSynced ..
	ErrAccountNotSynced = errors.New("account scan not synced to the highest")
	// ErrConnotGetBalance ..
	ErrConnotGetBalance = errors.New("cannot get balance from iron node")
	// ErrImportWalletWrong ..
	ErrImportWalletWrong = errors.New("import wallet failed")
	// ErrNotAcceptTx ..
	ErrNotAcceptTx = errors.New("node not accpect transaction")
	// ErrImportWalletWrong ..
	ErrTxNotSynced = errors.New("transaction have not be synced")
	// ErrImportWalletWrong ..
	ErrTransactionFailed = errors.New("ironfish transaction failed")
	// ErrImportWalletWrong ..
	ErrFeeToHigh = errors.New("fee is to high over the limit")
)

var (
	stopErrMsg    = []string{ErrNodeNotSynced.Error(), ErrAccountNotSynced.Error(), ErrTransactionFailed.Error(), ErrFeeToHigh.Error()}
	ironfishToken = &coins.TokenInfo{OfficialName: "IronFish", Decimal: 8, Unit: "IRON", Name: ChainNativeCoinName, OfficialContract: ChainNativeCoinName, TokenType: coins.Ironfish}
)

func init() {
	// set chain info
	ironfishToken.ChainType = ChainType
	ironfishToken.ChainNativeUnit = ChainNativeUnit
	ironfishToken.ChainAtomicUnit = ChainAtomicUnit
	ironfishToken.ChainUnitExp = ChainUnitExp
	ironfishToken.GasType = v1.GasType_GasUnsupported
	ironfishToken.ChainID = ChainID
	ironfishToken.ChainNickname = ChainNickname
	ironfishToken.ChainNativeCoinName = ChainNativeCoinName

	ironfishToken.Waight = 100
	ironfishToken.Net = coins.CoinNetMain
	ironfishToken.Contract = ironfishToken.OfficialContract
	ironfishToken.CoinType = sphinxplugin.CoinType_CoinTypeironfish
	register.RegisteTokenInfo(ironfishToken)
}

func ToIron(balance string) (decimal.Decimal, error) {
	//nolint
	balanceBigInt, ok := big.NewInt(0).SetString(balance, 10)
	if !ok {
		return decimal.Decimal{}, errors.New("balance is invalid")
	}
	return decimal.NewFromBigInt(balanceBigInt, IronExp), nil
}

func ToPoint(value float64) decimal.Decimal {
	return decimal.NewFromFloat(value).Mul(decimal.NewFromBigInt(big.NewInt(1), -IronExp))
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

package coins

import (
	"fmt"
	"strings"
	"time"

	v1 "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/NpoolPlatform/message/npool/sphinxplugin"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/utils"
)

type (
	TokenType string
)

const (
	Ironfish TokenType = "Ironfish"
)

type TokenInfo struct {
	OfficialName        string
	OfficialContract    string
	Contract            string // if ENV is main Contract = OfficialContract
	TokenType           TokenType
	Net                 string
	Unit                string
	Decimal             int
	Name                string
	Waight              int
	DisableRegiste      bool
	CoinType            sphinxplugin.CoinType
	ChainType           sphinxplugin.ChainType
	ChainNativeUnit     string
	ChainAtomicUnit     string
	ChainUnitExp        uint32
	GasType             v1.GasType
	ChainID             string
	ChainNickname       string
	ChainNativeCoinName string
}

const (
	CoinNetMain = "main"
	CoinNetTest = "test"
	TestPrefix  = "t"
)

var (
	// not export
	netCoinMap = map[string]map[string]sphinxplugin.CoinType{
		CoinNetMain: {
			"ironfish": sphinxplugin.CoinType_CoinTypeironfish,
		},
		CoinNetTest: {
			"ironfish": sphinxplugin.CoinType_CoinTypetironfish,
		},
	}

	// in order to compatible
	S3KeyPrxfixMap = map[string]string{
		"ironfish":  "ironfish/",
		"tironfish": "ironfish/",
	}

	// default sync time for waitting transaction on chain
	SyncTime = map[sphinxplugin.CoinType]time.Duration{
		sphinxplugin.CoinType_CoinTypeironfish:  time.Minute,
		sphinxplugin.CoinType_CoinTypetironfish: time.Minute,
	}
)

// CoinInfo report coin info
type CoinInfo struct {
	ENV      string // main or test
	Unit     string
	IP       string // wan ip
	Location string
}

// CheckSupportNet ..
func CheckSupportNet(netEnv string) bool {
	return (netEnv == CoinNetMain ||
		netEnv == CoinNetTest)
}

// TODO match case elegant deal
func CoinStr2CoinType(netEnv, coinStr string) sphinxplugin.CoinType {
	_netEnv := strings.ToLower(netEnv)
	_coinStr := strings.ToLower(coinStr)
	return netCoinMap[_netEnv][_coinStr]
}

func ToTestChainType(chainType sphinxplugin.ChainType) sphinxplugin.ChainType {
	if chainType == sphinxplugin.ChainType_UnKnow {
		return sphinxplugin.ChainType_UnKnow
	}
	name, ok := sphinxplugin.ChainType_name[int32(chainType)]
	if !ok {
		return sphinxplugin.ChainType_UnKnow
	}
	_chainType, ok := sphinxplugin.ChainType_value[fmt.Sprintf("T%v", name)]
	if !ok {
		return sphinxplugin.ChainType_UnKnow
	}
	return sphinxplugin.ChainType(_chainType)
}

func ToTestCoinType(coinType sphinxplugin.CoinType) sphinxplugin.CoinType {
	if coinType == sphinxplugin.CoinType_CoinTypeUnKnow {
		return sphinxplugin.CoinType_CoinTypeUnKnow
	}
	name := utils.ToCoinName(coinType)
	return CoinStr2CoinType(CoinNetTest, name)
}

func GetS3KeyPrxfix(tokenInfo *TokenInfo) string {
	if val, ok := S3KeyPrxfixMap[tokenInfo.Name]; ok {
		return val
	}

	name := tokenInfo.Name
	if tokenInfo.Net == CoinNetTest {
		name = strings.TrimPrefix(name, TestPrefix)
	}
	return fmt.Sprintf("%v/", name)
}

func GenerateName(tokenInfo *TokenInfo) string {
	chainType := utils.ToCoinName(tokenInfo.CoinType)
	name := strings.Trim(tokenInfo.OfficialName, " ")
	name = strings.ReplaceAll(name, " ", "-")
	return fmt.Sprintf("%v_%v_%v", chainType, tokenInfo.TokenType, name)
}

func GetChainType(in string) string {
	ret := strings.Split(in, "_")
	return ret[0]
}

package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/web3eye-io/ironfish-go-sdk/pkg/ironfish/api"
	"github.com/web3eye-io/ironfish-go-sdk/pkg/utils"

	//nolint
	"github.com/web3eye-io/ironfish-go-sdk/pkg/ironfish/types"

	"github.com/NpoolPlatform/message/npool/sphinxplugin"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/coins"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/coins/iron"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/coins/register"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/env"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/log"
	ct "github.com/NpoolPlatform/sphinx-plugin-p3/pkg/types"
)

// here register plugin func
func init() {
	register.RegisteTokenHandler(
		coins.Ironfish,
		register.OpGetBalance,
		walletBalance,
	)
	register.RegisteTokenHandler(
		coins.Ironfish,
		register.OpPreSign,
		preSign,
	)
	register.RegisteTokenHandler(
		coins.Ironfish,
		register.OpBroadcast,
		broadcast,
	)
	register.RegisteTokenHandler(
		coins.Ironfish,
		register.OpSyncTx,
		syncTx,
	)

	err := register.RegisteAbortFuncErr(sphinxplugin.CoinType_CoinTypeironfish, iron.TxFailErr)
	if err != nil {
		panic(err)
	}

	err = register.RegisteAbortFuncErr(sphinxplugin.CoinType_CoinTypetironfish, iron.TxFailErr)
	if err != nil {
		panic(err)
	}
}

func walletBalance(ctx context.Context, in []byte, tokenInfo *coins.TokenInfo) (out []byte, err error) {
	info := &iron.ViewAccount{}
	if err := json.Unmarshal(in, &info); err != nil {
		return in, err
	}

	v, ok := env.LookupEnv(env.ENVCOINNET)
	if !ok {
		return in, env.ErrEVNCoinNet
	}
	if !coins.CheckSupportNet(v) {
		return in, env.ErrEVNCoinNetValue
	}

	if info.PublicKey == "" {
		return in, env.ErrAddressInvalid
	}

	err = json.Unmarshal(in, info)
	if err != nil {
		return in, err
	}
	client := iron.Client()
	var bl *types.GetBalanceResponse
	err = client.WithClient(ctx, func(_ctx context.Context, cli *sdk.Client) (bool, error) {
		_, err := cli.ImportAccount(&types.ImportAccountRequest{
			Account: types.Account{
				Version:         info.Version,
				Name:            info.Name,
				PublicAddress:   info.PublicKey,
				ViewKey:         info.ViewKey,
				IncomingViewKey: info.IncomingKey,
				OutgoingViewKey: info.OutgoingKey,
			},
			Rescan: true,
		})

		if err != nil && !strings.Contains(err.Error(), "Account already exists") {
			return true, fmt.Errorf("%v ,%v", iron.ErrImportWalletWrong, err)
		}

		bl, err = cli.GetBalance(&types.GetBalanceRequest{
			Account:       info.Name,
			Confirmations: iron.DefaultConfirmations,
		})

		if err != nil {
			return true, err
		}
		if bl == nil {
			return true, iron.ErrConnotGetBalance
		}

		nodeStatus, err := cli.GetNodeStatus()
		if err != nil {
			return true, err
		}

		if nodeStatus.Blockchain.Head.Sequence-iron.ToleranceHeight > int(bl.Sequence) {
			return true, iron.ErrAccountNotSynced
		}
		return false, err
	})
	if err != nil {
		return in, err
	}

	balance, err := iron.ToIron(bl.Available)
	if err != nil {
		return in, fmt.Errorf("transform balance failed from %v", err)
	}
	f, exact := balance.Float64()
	if !exact {
		log.Warnf("wallet balance transfer warning balance %v from->to %v-%v", bl.Available, balance.String(), f)
	}

	_out := ct.WalletBalanceResponse{
		Balance:    f,
		BalanceStr: balance.String(),
	}
	return json.Marshal(_out)
}

func preSign(ctx context.Context, in []byte, tokenInfo *coins.TokenInfo) (out []byte, err error) {
	info := ct.BaseInfo{}
	if err := json.Unmarshal(in, &info); err != nil {
		return in, err
	}

	if !coins.CheckSupportNet(info.ENV) {
		return nil, env.ErrEVNCoinNetValue
	}

	amount := iron.ToPoint(info.Value).BigInt().Uint64()
	client := iron.Client()

	var createTxResp *types.CreateTransactionResponse
	err = client.WithClient(ctx, func(ctx context.Context, c *sdk.Client) (bool, error) {
		eFRResp, err := c.EstimateFeeRates()
		if err != nil {
			return true, err
		}
		feeRate := eFRResp.Average
		createTxResp, err = c.CreateTransaction(&types.CreateTransactionRequest{
			Account: info.From,
			Outputs: []types.Output{{
				PublicAddress: info.To,
				Amount:        fmt.Sprint(amount),
				Memo:          info.Memo,
			}},
			FeeRate: feeRate,
		})

		if err != nil {
			return true, err
		}

		fee, err := utils.GetFee(createTxResp.Transaction)
		if err != nil {
			return false, err
		}

		if fee > iron.MaxFeeLimit {
			return false, fmt.Errorf("%v, fee:%v, feeLimit: %v, feeRate: %v", iron.ErrFeeToHigh, fee, iron.MaxFeeLimit, feeRate)
		}

		return false, nil
	})

	if err != nil {
		return in, fmt.Errorf("%v,%v", iron.ErrTransactionFailed, err)
	}

	out, err = json.Marshal(iron.SignTxMsg{
		FromAccount: info.From,
		Transaction: createTxResp.Transaction,
	})
	if err != nil {
		return in, err
	}
	return out, nil
}

func broadcast(ctx context.Context, in []byte, tokenInfo *coins.TokenInfo) (out []byte, err error) {
	info := &iron.BroadcastTxMsg{}
	if err := json.Unmarshal(in, &info); err != nil {
		return in, err
	}

	client := iron.Client()

	var addTxResp *types.AddTransactionResponse
	err = client.WithClient(ctx, func(ctx context.Context, c *sdk.Client) (bool, error) {
		addTxResp, err = c.AddTransaction(&types.AddTransactionRequest{Transaction: info.SignedTransaction, Broadcast: true})
		return false, err
	})

	if err != nil {
		return in, fmt.Errorf("%v,%v", iron.ErrTransactionFailed, err)
	}
	out, err = json.Marshal(iron.SyncTxMsg{
		FromAccount: info.FromAccount,
		TxHash:      addTxResp.Hash,
	})
	if err != nil {
		return in, err
	}
	return out, nil
}

// syncTx sync transaction status on chain
func syncTx(ctx context.Context, in []byte, tokenInfo *coins.TokenInfo) (out []byte, err error) {
	info := &iron.SyncTxMsg{}
	if err := json.Unmarshal(in, &info); err != nil {
		return in, err
	}

	client := iron.Client()

	var getATxResp *types.GetAccountTransactionResponse
	err = client.WithClient(ctx, func(ctx context.Context, c *sdk.Client) (bool, error) {
		getATxResp, err = c.GetAccountTransaction(&types.GetAccountTransactionRequest{
			Hash:          info.TxHash,
			Account:       info.FromAccount,
			Confirmations: iron.DefaultConfirmations,
		})
		if err != nil {
			return true, err
		}
		return false, nil
	})

	if err != nil {
		return in, err
	}

	switch getATxResp.Transaction.Status {
	case types.EXPIRED:
		return in, iron.ErrTransactionFailed
	case types.CONFIRMED:
		return in, nil
	default:
		return in, iron.ErrTxNotSynced
	}
}

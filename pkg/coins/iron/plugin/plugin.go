package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/web3eye-io/ironfish-go-sdk/pkg/ironfish/api"
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
			Confirmations: 2,
		})

		if err != nil {
			return true, err
		}
		if bl == nil {
			return true, iron.ErrConnotGetBalance
		}

		nodeStatus, err := cli.GetNodeStatus()
		// TODO: will be confirmed ,how sequence running
		if nodeStatus.Blockchain.Head.Sequence-10 > int(bl.Sequence) {
			return true, iron.ErrAccountNotSynced
		}
		return false, err
	})
	if err != nil {
		return in, err
	}

	bigBalance, ok := big.NewInt(0).SetString(bl.Available, 10)
	if !ok {
		return in, fmt.Errorf("transform balance failed from %v", bl.Available)
	}
	balance := iron.ToIron(bigBalance.Uint64())
	f, exact := balance.Float64()
	if exact != big.Exact {
		log.Warnf("wallet balance transfer warning balance from->to %v-%v", balance.String(), f)
	}

	_out := ct.WalletBalanceResponse{
		Balance:    f,
		BalanceStr: balance.String(),
	}
	return json.Marshal(_out)
}

func preSign(ctx context.Context, in []byte, tokenInfo *coins.TokenInfo) (out []byte, err error) {
	// info := ct.BaseInfo{}
	// if err := json.Unmarshal(in, &info); err != nil {
	// 	return in, err
	// }

	// if !coins.CheckSupportNet(info.ENV) {
	// 	return nil, env.ErrEVNCoinNetValue
	// }

	// client := sol.Client()

	// var recentBlockHash *rpc.GetLatestBlockhashResult
	// err = client.WithClient(ctx, func(_ctx context.Context, cli *rpc.Client) (bool, error) {
	// 	recentBlockHash, err = cli.GetLatestBlockhash(_ctx, rpc.CommitmentFinalized)
	// 	if err != nil || recentBlockHash == nil {
	// 		return true, err
	// 	}
	// 	return false, err
	// })
	// if err != nil {
	// 	return in, err
	// }

	// _out := sol.SignMsgTx{
	// 	BaseInfo:        info,
	// 	RecentBlockHash: recentBlockHash.Value.Blockhash.String(),
	// }

	// return json.Marshal(_out)
	return nil, nil

}

func broadcast(ctx context.Context, in []byte, tokenInfo *coins.TokenInfo) (out []byte, err error) {
	// info := sol.BroadcastRequest{}
	// if err := json.Unmarshal(in, &info); err != nil {
	// 	return in, err
	// }

	// tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(info.Signature))
	// if err != nil {
	// 	return in, err
	// }

	// err = tx.VerifySignatures()
	// if err != nil {
	// 	return in, sol.ErrSolSignatureWrong
	// }

	// client := sol.Client()
	// if err != nil {
	// 	return in, err
	// }
	// var cid solana.Signature
	// err = client.WithClient(ctx, func(_ctx context.Context, cli *rpc.Client) (bool, error) {
	// 	cid, err = cli.SendTransaction(_ctx, tx)
	// 	if err != nil && !sol.TxFailErr(err) {
	// 		return true, err
	// 	}
	// 	return false, err
	// })
	// if err != nil {
	// 	return in, err
	// }

	// _out := ct.SyncRequest{
	// 	TxID: cid.String(),
	// }

	// return json.Marshal(_out)
	return nil, nil
}

// syncTx sync transaction status on chain
func syncTx(ctx context.Context, in []byte, tokenInfo *coins.TokenInfo) (out []byte, err error) {
	// info := ct.SyncRequest{}
	// if err := json.Unmarshal(in, &info); err != nil {
	// 	return in, err
	// }

	// signature, err := solana.SignatureFromBase58(info.TxID)
	// if err != nil {
	// 	return in, err
	// }

	// client := sol.Client()
	// var chainMsg *rpc.GetTransactionResult
	// err = client.WithClient(ctx, func(_ctx context.Context, cli *rpc.Client) (bool, error) {
	// 	chainMsg, err = cli.GetTransaction(
	// 		_ctx,
	// 		signature,
	// 		&rpc.GetTransactionOpts{
	// 			Encoding:   solana.EncodingBase58,
	// 			Commitment: rpc.CommitmentFinalized,
	// 		})
	// 	if err != nil {
	// 		return true, err
	// 	}
	// 	return false, err
	// })

	// if err != nil {
	// 	return in, err
	// }

	// if chainMsg == nil {
	// 	return in, env.ErrWaitMessageOnChain
	// }

	// if chainMsg != nil && chainMsg.Meta.Err != nil {
	// 	sResp := &ct.SyncResponse{}
	// 	sResp.ExitCode = -1
	// 	out, mErr := json.Marshal(sResp)
	// 	if mErr != nil {
	// 		return in, mErr
	// 	}
	// 	return out, fmt.Errorf("%v,%v", sol.SolTransactionFailed, err)
	// }

	// if chainMsg != nil && chainMsg.Meta.Err == nil {
	// 	sResp := &ct.SyncResponse{}
	// 	sResp.ExitCode = 0
	// 	out, err := json.Marshal(sResp)
	// 	if err != nil {
	// 		return in, err
	// 	}
	// 	return out, nil
	// }

	// return in, sol.ErrSolBlockNotFound
	return nil, nil

}

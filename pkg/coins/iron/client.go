package iron

import (
	"context"
	"strings"
	"time"

	sdk "github.com/web3eye-io/ironfish-go-sdk/pkg/ironfish/api"
	//nolint
	"github.com/web3eye-io/ironfish-go-sdk/pkg/ironfish/types"

	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/endpoints"
	"github.com/NpoolPlatform/sphinx-plugin-p3/pkg/utils"
)

const (
	MinNodeNum       = 1
	MaxRetries       = 3
	retriesSleepTime = 200 * time.Millisecond
	reqTimeout       = 6 * time.Second
	connectTimeout   = 3 * time.Second
)

type IRClientI interface {
	GetNode(ctx context.Context, endpointmgr *endpoints.Manager) (*sdk.Client, error)
	WithClient(ctx context.Context, fn func(context.Context, *sdk.Client) (bool, error)) error
}

type IRClients struct{}

func (irClients IRClients) GetNode(_ctx context.Context, endpointmgr *endpoints.Manager) (*sdk.Client, error) {
	endpoint, err := endpointmgr.Peek()
	if err != nil {
		return nil, err
	}

	_, cancel := context.WithTimeout(_ctx, reqTimeout)
	defer cancel()
	addr, authToken := "", ""
	segStr := strings.Split(endpoint, "|")
	addr = segStr[0]
	if len(segStr) < 2 {
		authToken = ""
	} else {
		authToken = segStr[1]
	}
	client := sdk.NewClient(addr, authToken, true)
	err = client.Connect(connectTimeout)
	if err != nil {
		return nil, err
	}

	nodeStatus, err := client.GetNodeStatus()
	if err != nil {
		return nil, err
	}

	if nodeStatus.Node.Status != types.NodeStarted &&
		nodeStatus.BlockSyncer.Status != types.BlockSyncerSyncing &&
		!nodeStatus.Blockchain.Synced ||
		!nodeStatus.PeerNetwork.IsReady {
		return nil, ErrNodeNotSynced
	}

	return client, nil
}

func (irClients *IRClients) WithClient(ctx context.Context, fn func(ctx context.Context, c *sdk.Client) (bool, error)) error {
	var (
		apiErr, err error
		retry       bool
		client      *sdk.Client
	)
	endpointmgr, err := endpoints.NewManager()
	if err != nil {
		return err
	}

	for i := 0; i < utils.MinInt(MaxRetries, endpointmgr.Len()); i++ {
		if i > 0 {
			time.Sleep(retriesSleepTime)
		}

		client, err = irClients.GetNode(ctx, endpointmgr)
		if err != nil {
			continue
		}

		retry, apiErr = fn(ctx, client)
		if !retry {
			return apiErr
		}
	}
	if apiErr != nil {
		return apiErr
	}
	return err
}

func Client() IRClientI {
	return &IRClients{}
}

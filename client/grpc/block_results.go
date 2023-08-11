package grpc

import (
	"context"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/gogoproto/grpc"

	brs "github.com/comet/comet-companion/client/proto/v1"
)

type BlockResults struct {
	Height                int64                     `json:"height"`
	TxsResults            []*abci.ExecTxResult      `json:"txs_results"`
	FinalizeBlockEvents   []*abci.Event             `json:"finalize_block_events"`
	ValidatorUpdates      []*abci.ValidatorUpdate   `json:"validator_updates"`
	ConsensusParamUpdates *cmtproto.ConsensusParams `json:"consensus_param_updates"`
	AppHash               []byte                    `json:"app_hash"`
}

// BlockResultsServiceClient provides the block results of a given height (or latest if none provided).
type BlockResultsServiceClient interface {
	GetBlockResults(ctx context.Context, height int64) (*BlockResults, error)
	GetLatestBlockResults(ctx context.Context) (*BlockResults, error)
}

type blockResultServiceClient struct {
	client brs.BlockResultsServiceClient
}

func (b blockResultServiceClient) GetBlockResults(ctx context.Context, height int64) (*BlockResults, error) {
	res, err := b.client.GetBlockResults(ctx, &brs.GetBlockResultsRequest{Height: height})
	if err != nil {
		return nil, fmt.Errorf("error fetching BlockResults for height %d:: %s", height, err.Error())
	}

	return &BlockResults{
		Height:                res.Height,
		TxsResults:            res.TxsResults,
		FinalizeBlockEvents:   res.FinalizeBlockEvents,
		ValidatorUpdates:      res.ValidatorUpdates,
		ConsensusParamUpdates: res.ConsensusParamUpdates,
		AppHash:               res.AppHash,
	}, nil
}

func (b blockResultServiceClient) GetLatestBlockResults(ctx context.Context) (*BlockResults, error) {
	res, err := b.client.GetLatestBlockResults(ctx, &brs.GetLatestBlockResultsRequest{})
	if err != nil {
		return nil, fmt.Errorf("error fetching BlockResults for latest height :: %s", err.Error())
	}

	return &BlockResults{
		Height:                res.Height,
		TxsResults:            res.TxsResults,
		FinalizeBlockEvents:   res.FinalizeBlockEvents,
		ValidatorUpdates:      res.ValidatorUpdates,
		ConsensusParamUpdates: res.ConsensusParamUpdates,
		AppHash:               res.AppHash,
	}, nil
}

func newBlockResultsServiceClient(conn grpc.ClientConn) BlockResultsServiceClient {
	return &blockResultServiceClient{
		client: brs.NewBlockResultsServiceClient(conn),
	}
}

type disabledBlockResultsServiceClient struct{}

func newDisabledBlockResultsServiceClient() BlockResultsServiceClient {
	return &disabledBlockResultsServiceClient{}
}

// GetBlockResults implements BlockResultsServiceClient
func (*disabledBlockResultsServiceClient) GetBlockResults(_ context.Context, _ int64) (*BlockResults, error) {
	panic("block results service client is disabled")
}

// GetLatestBlockResults implements BlockResultsServiceClient
func (*disabledBlockResultsServiceClient) GetLatestBlockResults(_ context.Context) (*BlockResults, error) {
	panic("block results service client is disabled")
}
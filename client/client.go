package client

import (
	"context"
	"fmt"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	cometClient "github.com/cometbft/cometbft/rpc/grpc/client"
	"github.com/spf13/viper"
	"time"
)

func CreateRPCConnection() *rpchttp.HTTP {
	addr := viper.GetString("app.rpc")
	conn, err := rpchttp.New(addr, "/websocket")
	if err != nil {
		fmt.Errorf("failed to dial %s: %w", addr, err)
	}
	return conn
}

func CreateGRPCConnection(ctx context.Context) cometClient.Client {
	addr := viper.GetString("app.grpc")

	conn, err := cometClient.New(ctx, addr, cometClient.WithInsecure())
	if err != nil {
		fmt.Errorf("failed to dial %s: %w", addr, err)
	}
	return conn
}

func ProcessReq(request string, latest bool, height int) {
	if request == "Block" {
		GetBlock(latest, int64(height))
	} else if request == "BlockResult" {
		GetBlockResults(latest, int64(height))
	} else {
		fmt.Printf("Error Invalid Input: %s, %v", request, latest)
	}
}

func GetBlock(latest bool, height int64) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer ctxCancel()
	conn := CreateGRPCConnection(ctx)
	var res *cometClient.Block
	var err error

	if latest {
		res, err = conn.GetLatestBlock(ctx)
	} else {
		res, err = conn.GetBlockByHeight(ctx, height)
	}
	if err != nil {
		fmt.Printf("Error making request: %w", err)
	}
	fmt.Printf("Block ID: %s\nBlock: %s", res.BlockID.String(), res.Block.String())
}

func GetBlockResults(latest bool, height int64) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer ctxCancel()
	conn := CreateGRPCConnection(ctx)
	var res *cometClient.BlockResults
	var err error
	if latest {
		res, err = conn.GetLatestBlockResults(ctx)
	} else {
		res, err = conn.GetBlockResults(ctx, height)
	}
	if err != nil {
		fmt.Printf("Error making request: %w", err)
	}
	fmt.Printf("Height: %v\nTxResults: %v\nFinalizeBlockEvents: %s\nValidatorUpdates: %s\nConsensusParamUpdatesL %s\nAppHash: %s\n",
		res.Height, res.TxsResults, res.FinalizeBlockEvents, res.ValidatorUpdates, res.ConsensusParamUpdates, string(res.AppHash))
}

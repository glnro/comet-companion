package main

import (
	"context"
	"fmt"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	cometClient "github.com/cometbft/cometbft/rpc/grpc/client"
	"time"
)

func main() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer ctxCancel()

	benchmarkGRPC(ctx)
	benchmarkRPC(ctx)
}

func benchmarkGRPC(ctx context.Context) {
	conn := createGRPCConnection(ctx)
	reqStart := time.Now()

	getGRPCResult(func() (any, error) { return conn.GetVersion(ctx) }, "GRPC :: Get Version response took %v ms")
	getGRPCResult(func() (any, error) { return conn.GetBlockResults(ctx, 1) }, "GRPC :: Get Block Results took %v ms")
	getGRPCResult(func() (any, error) { return conn.GetLatestBlockResults(ctx) }, "GRPC :: Get Latest Block Results took %v ms")
	getGRPCResult(func() (any, error) { return conn.GetLatestBlock(ctx) }, "GRPC :: Get Latest Block took %v ms")
	getGRPCResult(func() (any, error) { return conn.GetBlockByHeight(ctx, 1) }, "GRPC :: Get Block took %v ms")

	reqEnd := time.Now()
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	fmt.Println(fmt.Sprintf("GRPC requests took %v ms", elapsed))
}

func benchmarkRPC(ctx context.Context) {
	conn := createRPCConnection()
	height := int64(1)
	reqStart := time.Now()
	getRPCResult(func() (any, error) { return conn.BlockResults(ctx, &height) }, "RPC :: Get Block Results took %v ms")
	getRPCResult(func() (any, error) { return conn.BlockResults(ctx, nil) }, "RPC :: Get Latest Block Results took %v ms")
	getRPCResult(func() (any, error) { return conn.Block(ctx, nil) }, "RPC :: Get Latest Block took %v ms")
	getRPCResult(func() (any, error) { return conn.Block(ctx, &height) }, "RPC :: Get Block took %v ms")
	reqEnd := time.Now()
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	fmt.Println(fmt.Sprintf("RPC requests took %v ms", elapsed))
}

func createRPCConnection() *rpchttp.HTTP {
	addr := "http://127.0.0.1:5701"
	conn, err := rpchttp.New(addr, "/websocket")
	if err != nil {
		fmt.Errorf("failed to dial %s: %w", addr, err)
	}
	return conn
}

func createGRPCConnection(ctx context.Context) cometClient.Client {
	addr := "127.0.0.1:5702"
	conn, err := cometClient.New(ctx, addr, cometClient.WithInsecure())
	if err != nil {
		fmt.Errorf("failed to dial %s: %w", addr, err)
	}
	return conn
}

func getGRPCResult(grpcCall func() (any, error), errMsg string) {
	reqStart := time.Now()
	_, err := grpcCall()
	reqEnd := time.Now()
	if err != nil {
		fmt.Errorf("failed to retrieve response: %w", err)
	}
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	fmt.Println(fmt.Sprintf(errMsg, elapsed))
}

func getRPCResult(httpCall func() (any, error), errMsg string) {
	reqStart := time.Now()
	_, err := httpCall()
	reqEnd := time.Now()
	if err != nil {
		fmt.Errorf("failed to retrieve response: %w", err)
	}
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	fmt.Println(fmt.Sprintf(errMsg, elapsed))
}

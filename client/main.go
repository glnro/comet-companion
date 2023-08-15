package main

import (
	"context"
	"fmt"
	cl "github.com/comet/comet-companion/client/client/grpc"
	cometClient "github.com/cometbft/cometbft/rpc/grpc/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func main() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer ctxCancel()
	conn := createGRPCConnection(ctx)

	getVersion(ctx, conn)
	getBlockResults(ctx, conn)
	getLatestBlockResults(ctx, conn)
	getLatestBlock(ctx, conn)
	getBlock(ctx, conn)
}

func createGRPCConnection(ctx context.Context) cometClient.Client {
	addr := "127.0.0.1:5702"

	conn, err := cometClient.New(ctx, addr, cometClient.WithInsecure())
	if err != nil {
		fmt.Errorf("failed to dial %s: %w", addr, err)
	}
	return conn
}

func getVersion(ctx context.Context, conn cometClient.Client) {
	reqStart := time.Now()
	_, err := conn.GetVersion(ctx)
	reqEnd := time.Now()
	if err != nil {
		fmt.Errorf("failed to retrieve version: %w", err)
	}
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	fmt.Println(fmt.Sprintf("Get Version response took %v ms", elapsed))
}

func getBlockResults(ctx context.Context, conn cometClient.Client) {
	reqStart := time.Now()
	_, err := conn.GetBlockResults(ctx, 1)
	reqEnd := time.Now()
	if err != nil {
		fmt.Errorf("failed to retrieve block Results: %w", err)
	}
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	fmt.Println(fmt.Sprintf("Get Block Results took %v ms", elapsed))
}

func getLatestBlockResults(ctx context.Context, conn cometClient.Client) {
	reqStart := time.Now()
	_, err := conn.GetLatestBlockResults(ctx)
	reqEnd := time.Now()
	if err != nil {
		fmt.Errorf("failed to retrieve latest block Results: %w", err)
	}
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	fmt.Println(fmt.Sprintf("Get Latest Block Results took %v ms", elapsed))
}

func getLatestBlock(ctx context.Context, conn cometClient.Client) {
	reqStart := time.Now()
	_, err := conn.GetLatestBlock(ctx)
	reqEnd := time.Now()
	if err != nil {
		fmt.Errorf("failed to retrieve block Results: %w", err)
	}
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	fmt.Println(fmt.Sprintf("Get Latest Block took %v ms", elapsed))
}

func getBlock(ctx context.Context, conn cometClient.Client) {
	reqStart := time.Now()
	_, err := conn.GetBlockByHeight(ctx, 1)
	reqEnd := time.Now()
	if err != nil {
		fmt.Errorf("failed to retrieve block Results: %w", err)
	}
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	fmt.Println(fmt.Sprintf("Get Block took %v ms", elapsed))
}

func vanillaGrpc() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer ctxCancel()
	addr := "127.0.0.1:5702"

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("failed to dial %s: %w", addr, err)
	}
	defer conn.Close()
	fmt.Println("connected to client")

	versionClient := cl.NewVersionServiceClient(conn)

	res, err := versionClient.GetVersion(ctx)
	if err != nil {
		fmt.Errorf("failed to retrieve version: %s: %w", addr, err)
	}
	fmt.Println(fmt.Sprintf("Response: %v", res.ABCI))
}

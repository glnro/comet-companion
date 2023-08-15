package main

import (
	"context"
	"fmt"
	cl "github.com/comet/comet-companion/client/client/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func main() {
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

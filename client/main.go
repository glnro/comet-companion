package main

import (
	"fmt"
	"log"

	"context"
	cometclient "github.com/cometbft/cometbft/rpc/grpc/client"
	"time"
)

func main() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer ctxCancel()

	client, err := cometclient.New(
		ctx,
		fmt.Sprintf("localhost:%v", 5702),
		cometclient.WithInsecure(),
	)
	if err != nil {
		log.Panicf("Unable to connect to grpc")
	}

	fmt.Println("Connected to client")

	for {
		res, err := client.GetVersion(ctx)
		if err != nil {
			log.Fatalf("Unable to fetch latest client %v", err.Error())
		}

		fmt.Sprintf("Response: %v", res)
	}
}

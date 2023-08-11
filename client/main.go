package main

import (
	"fmt"
	"log"

	"context"
	grpcclient "github.com/comet/comet-companion/client/client/grpc"
	"time"
)

func main() {
	fmt.Println("hello world")

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Minute)
	defer ctxCancel()

	client, err := grpcclient.New(
		ctx,
		fmt.Sprintf("localhost:%v", 5702),
		grpcclient.WithInsecure(),
	)
	if err != nil {
		log.Panicf("Unable to connect to grpc")
	}

	fmt.Println("Connected to client")

	// Wait for comet to produce a few blocks
	//time.Sleep(time.Second * 60)

	for {
		res, err := client.GetVersion(ctx)
		if err != nil {
			log.Fatalf("Unable to fetch latest client %v", err.Error())
		}

		fmt.Sprintf("Response: %v", res)
	}
}

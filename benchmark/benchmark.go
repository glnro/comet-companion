package benchmark

import (
	"context"
	"flag"
	"fmt"
	"github.com/comet/comet-companion/client/client"
	"github.com/cometbft/cometbft/libs/rand"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	cometClient "github.com/cometbft/cometbft/rpc/grpc/client"
	"math"
	"sync"
	"time"
)

type Result struct {
	Low  int64
	Avg  int64
	High int64
}

func benchmark() {
	requestVol := *flag.Int("reqVol", 10, "Number of requests to send")
	request := *flag.String("request", "LatestBlock", "Endpoint to benchmark")
	grpc := *flag.Bool("GRPC", true, "Test GRPC (RPC=false")
	flag.Parse()
	runTest(requestVol, request, grpc)
}

func runTest(requestVol int, request string, grpcOpt bool) {
	finalResults := []Result{}
	n := requestVol
	for i := 0; i < 10; i++ {
		ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*30)
		defer ctxCancel()

		results := make(chan int64, n)
		runner(ctx, grpcOpt, results, n, request)
		finalResults = append(finalResults, getAverage(results, int64(n)))
	}
	final := getFinal(finalResults)
	fmt.Printf("Averages for 10x %s runs at %d\n", request, requestVol)
	fmt.Printf("Minimum Response Time: %v\n", final.Low)
	fmt.Printf("Average Response Time: %v\n", final.Avg)
	fmt.Printf("Maximum Response Time: %v\n", final.High)
}

func runner(ctx context.Context, grpc bool, results chan<- int64, n int, request string) {
	if grpc {
		conn := client.CreateGRPCConnection(ctx)
		defer conn.Close()
		height, err := conn.GetLatestBlock(ctx)
		if err != nil {
			fmt.Errorf("Error fetching latest block: %w", err)
		}
		h := rand.Intn(int(height.Block.Height))
		for i := 0; i < n; i++ {
			switch request {
			case "LatestBlock":
				reportGRPCResult(func() (any, error) { return conn.GetLatestBlock(ctx) }, results)
			case "BlockHeight":
				reportGRPCResult(func() (any, error) { return conn.GetBlockByHeight(ctx, int64(h)) }, results)
			case "LatestBlockResults":
				reportGRPCResult(func() (any, error) { return conn.GetLatestBlockResults(ctx) }, results)
			case "BlockResultsHeight":
				reportGRPCResult(func() (any, error) { return conn.GetBlockResults(ctx, int64(h)) }, results)
			default:
				return
			}

		}
		close(results)
	} else {
		conn := client.CreateRPCConnection()
		height, err := conn.Block(ctx, nil)
		if err != nil {
			fmt.Errorf("Error fetching latest block: %w", err)
		}
		h := int64(rand.Intn(int(height.Block.Height)))
		for i := 0; i < n; i++ {
			switch request {
			case "LatestBlock":
				reportGRPCResult(func() (any, error) { return conn.Block(ctx, nil) }, results)
			case "BlockHeight":
				reportGRPCResult(func() (any, error) { return conn.Block(ctx, &h) }, results)
			case "LatestBlockResults":
				reportGRPCResult(func() (any, error) { return conn.BlockResults(ctx, nil) }, results)
			case "BlockResultsHeight":
				reportGRPCResult(func() (any, error) { return conn.BlockResults(ctx, &h) }, results)
			default:
				panic("Invalid Option")
			}
		}
		close(results)
	}
}

func getFinal(results []Result) Result {
	var (
		totalLow  int64
		totalAvg  int64
		totalHigh int64
	)
	for _, r := range results {
		totalLow += r.Low
		totalAvg += r.Avg
		totalHigh += r.High
	}
	totalLen := int64(len(results))
	return Result{
		Low:  totalLow / totalLen,
		Avg:  totalAvg / totalLen,
		High: totalHigh / totalLen,
	}
}

func getAverage(results <-chan int64, n int64) Result {
	var (
		totalTime   int64
		minDuration int64 = math.MaxInt64
		maxDuration int64
		count       int
	)

	for result := range results {
		totalTime += result
		count++
		if result < minDuration {
			minDuration = result
		}
		if result > maxDuration {
			maxDuration = result
		}
	}

	averageDuration := totalTime / n

	return Result{
		Low:  minDuration,
		Avg:  averageDuration,
		High: maxDuration,
	}
}

func concurrentGrpcResult(conn cometClient.Client, res chan<- int64, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer ctxCancel()

	reqStart := time.Now()
	_, err := conn.GetLatestBlock(ctx)
	reqEnd := time.Now()
	if err != nil {
		fmt.Errorf("failed to retrieve response: %w", err)
	}
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	res <- elapsed
}

func concurrentRpcResult(conn *rpchttp.HTTP, res chan<- int64, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer ctxCancel()

	reqStart := time.Now()
	_, err := conn.Block(ctx, nil)
	reqEnd := time.Now()
	if err != nil {
		fmt.Errorf("failed to retrieve response: %w", err)
	}
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	res <- elapsed
}

func reportGRPCResult(grpcCall func() (any, error), res chan<- int64) {
	reqStart := time.Now()
	_, err := grpcCall()
	reqEnd := time.Now()
	if err != nil {
		fmt.Errorf("failed to retrieve response: %w", err)
	}
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	res <- elapsed
}

func reportRPCResult(httpCall func() (any, error), res chan<- int64) {
	reqStart := time.Now()
	_, err := httpCall()
	reqEnd := time.Now()
	if err != nil {
		fmt.Errorf("failed to retrieve response: %w", err)
	}
	elapsed := reqEnd.Sub(reqStart).Milliseconds()
	res <- elapsed
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

func BenchmarkGRPCPrint(ctx context.Context) {
	conn := client.CreateGRPCConnection(ctx)
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

func BenchmarkRPCPPrint(ctx context.Context) {
	conn := client.CreateRPCConnection()
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

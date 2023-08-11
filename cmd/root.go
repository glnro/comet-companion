package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/comet/comet-companion/client/benchmark"
	"github.com/comet/comet-companion/client/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"time"
)

const (
	Block       = "Block"
	BlockResult = "BlockResult"
	GRPC        = "grpc"
	RPC         = "rpc"
)

var rootCmd = &cobra.Command{
	Use:   "companion",
	Short: "Comet Companion App",
	Long: `Comet Companion interfaces 
   
Benchmark or perform requests against Comet's GRPC and RPC endpoints'`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Unexpected Error: '%s'", err)
		os.Exit(1)
	}
}

func init() {
	viper.SetConfigFile("../config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("Error reading config %w", err)
	}

	fmt.Printf("GRPC: %s", viper.GetString("app.grpc"))

	rootCmd.AddCommand(NewBenchmarkCommand())
	rootCmd.AddCommand(NewGrpcCommand())
}

func NewBenchmarkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "benchmark [grpc || rpc]",
		Short: "Run benchmark tests",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*30)
			defer ctxCancel()
			requestType := args[0]
			if requestType == GRPC {
				benchmark.BenchmarkGRPCPrint(ctx)
			} else if requestType == RPC {
				benchmark.BenchmarkRPCPPrint(ctx)
			} else {
				errors.New("Invalid request type [grpc || rpc]")
			}
			return nil
		},
	}
	return cmd
}

func NewGrpcCommand() *cobra.Command {
	grpcCmd := &cobra.Command{
		Use:                        "grpc",
		Short:                      "GRPC client subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	grpcCmd.AddCommand(NewBlockCommand())
	grpcCmd.AddCommand(NewBlockResultCommand())
	return grpcCmd
}

func NewBlockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block [latest] --height",
		Short: "Get latest or specified block height",
		Long: `Get Block
   
To get the latest height, set [latest] to false
To get a specific height, set [latest] to true and --height to desired block height to retrieve.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			latest, err := strconv.ParseBool(args[0])
			if err != nil {
				return fmt.Errorf("Invalid arg [latest]: %w", err)
			}
			height, err := cmd.Flags().GetInt("height")
			if !latest && height == 0 {
				return fmt.Errorf("Invalid height for Get Block latest=false")
			}
			client.ProcessReq(Block, latest, height)

			return nil
		},
	}
	cmd.Flags().Int("height", 0, "Desired Block height to retrieve")
	return cmd
}

func NewBlockResultCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block-result [latest] --height",
		Short: "Get latest or specified block height",
		Long: `Get BlockResult
   
To get the latest height, set [latest] to false
To get a specific height, set [latest] to true and --height to desired height to retrieve.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			latest, err := strconv.ParseBool(args[0])
			if err != nil {
				return fmt.Errorf("Invalid arg [latest]: %w", err)
			}
			height, err := cmd.Flags().GetInt("height")
			if !latest && height == 0 {
				return fmt.Errorf("Invalid height for Get BlockResult latest=false")
			}
			client.ProcessReq(BlockResult, latest, height)

			return nil
		},
	}
	cmd.Flags().Int("height", 0, "Desired BlockResult height to retrieve")
	return cmd
}

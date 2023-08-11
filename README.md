## Comet Companion Client CLI

Comet Companion Client CLI offers two distinct commands, one for performing individual requests and one for benchmarking request peformance.

### Getting Started
1. Configure the CLI to point to the correct network endpoints in the config.yml
```yaml
# Default GRPC port is 26670 and RPC is 26657
app:
  gprc: "<address>:port"
  rpc: "http://<address>:<port>"
```
2. Build or install
```shell
go build -o companion
```

### GRPC Requests
The requests subcommand offers GRPC `Block` and `BlockResults` for the latest or a given height.

To request the latest height for both:
`companion grpc [latest] --height`
```shell
./companion grpc true
```
To request a specific height, pass false to latest and the height to the flag
```shell
./companion grpc false --height=10000
```

### Benchmarking

To run request benchmarks, run the benchmark subcommand. Future iterations will include more advanced benchmarking configuration.
`companion benchmark [grpc || grpc]`
```protobuf
// For GRPC
./companion benchmark grpc

// For RPC
./companion benchmark rpc
```
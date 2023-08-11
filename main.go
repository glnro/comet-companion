package main

import (
	"fmt"
	"github.com/comet/comet-companion/client/cmd"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("Error reading config %w\n", err)
	}

	cmd.Execute()
}

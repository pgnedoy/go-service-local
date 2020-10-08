package main

import (
	"github.com/spf13/viper"

	"github.com/pgnedoy/go-service-local/hello-grpc/cmd"
)

func init() {
	viper.AutomaticEnv()
}

func main() {
	cmd.Execute()
}

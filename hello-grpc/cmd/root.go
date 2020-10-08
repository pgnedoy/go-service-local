package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use: "hello-grpc",
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatal("error in rootCommand.Execute", err)
	}
}

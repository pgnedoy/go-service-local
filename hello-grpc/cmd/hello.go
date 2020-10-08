package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	hellogrpcv1 "github.com/pgnedoy/protos/gen/go/hellogrpc/v1"

	"github.com/pgnedoy/go-service-local/core/flags"
	"github.com/pgnedoy/go-service-local/core/grpc"
)

var helloFlags *flags.Flags

var helloCommand = &cobra.Command{
	Use:  "hello <name>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := helloFlags.GetString("endpoint")

		if endpoint == "" {
			log.Fatal("invalid: endpoint")
		}

		conn, err := grpc.NewServiceMeshConnection(
			endpoint,
			helloFlags.GetString("service-name"),
		)

		if err != nil {
			log.Fatal("could not establish connection", err)
		}

		defer func() {
			if closeErr := conn.Close(); closeErr != nil {
				log.Println("error while closing connection", closeErr)
			}
		}()

		c := hellogrpcv1.NewHelloAPIClient(conn)

		req := &hellogrpcv1.HelloRequest{
			Name: args[0],
		}

		resp, err := c.Hello(context.Background(), req)

		if err != nil {
			log.Fatal("error with request", err)
		}

		log.Println(fmt.Sprintf("%v", resp))
	},
}

func init() {
	rootCommand.AddCommand(helloCommand)

	helloFlags = flags.New("hello", helloCommand)

	helloFlags.RegisterString("endpoint", "e", "localhost:5000", "Server endpoint", "")
	helloFlags.RegisterString("service-name", "", "hello-grpc-local", "Name of service in the mesh", "")
}

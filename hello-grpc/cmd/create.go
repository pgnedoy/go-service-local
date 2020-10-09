package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	creategrpcv1 "github.com/pgnedoy/protos/gen/go/hello-proto/v1"

	"github.com/pgnedoy/go-service-local/core/flags"
	"github.com/pgnedoy/go-service-local/core/grpc"
)

var createFlags *flags.Flags

var createCommand = &cobra.Command{
	Use:  "create <name>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := createFlags.GetString("endpoint")

		if endpoint == "" {
			log.Fatal("invalid: endpoint")
		}

		conn, err := grpc.NewServiceMeshConnection(
			endpoint,
			createFlags.GetString("service-name"),
		)

		if err != nil {
			log.Fatal("could not establish connection", err)
		}

		defer func() {
			if closeErr := conn.Close(); closeErr != nil {
				log.Println("error while closing connection", closeErr)
			}
		}()

		c := creategrpcv1.NewHelloAPIClient(conn)

		req := &creategrpcv1.CreateUserRequest{
			Name: args[0],
			AuthType: creategrpcv1.AuthType_AUTH_TYPE_SNAP,
			AuthId: "123-df-213-sdf",
			Age: 32,
			Country: "UK",
		}

		resp, err := c.CreateUser(context.Background(), req)

		if err != nil {
			log.Fatal("error with request", err)
		}

		log.Println(fmt.Sprintf("%v", resp))
	},
}

func init() {
	rootCommand.AddCommand(createCommand)

	createFlags = flags.New("create", createCommand)

	createFlags.RegisterString("endpoint", "e", "localhost:5000", "Server endpoint", "")
	createFlags.RegisterString("service-name", "", "create-grpc-local", "Name of service in the mesh", "")
}

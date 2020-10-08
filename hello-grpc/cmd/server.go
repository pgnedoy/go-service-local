package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/reflection"

	hellogrpcv1 "github.com/pgnedoy/protos/gen/go/boilerplates/hellogrpc/v1"

	"github.com/pgnedoy/go-service-local/core/flags"
	"github.com/pgnedoy/go-service-local/core/grpc"
)

const (
	defaultGRPCPort = 5000
)

var serverFlags *flags.Flags

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Run the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			log.Println("cancelling root context")
			cancel()
		}()

		port := serverFlags.GetInt("port")

		log.Println(fmt.Sprintf("server listening on port: %d", port))

		s, err := grpc.NewServer(&grpc.ServerConfig{
			ServiceName: serverFlags.GetString("service-name"),
			Environment: serverFlags.GetString("environment"),
		})

		if err != nil {
			log.Fatal("error initializing grpc server", err)
		}

		helloServer, err := handlers.NewGrpcServer(&handlers.GrpcServerConfig{})

		if err != nil {
			log.Fatal("error with NewGrpcServer", err)
		}

		hellogrpcv1.RegisterHelloAPIServer(s, helloServer)

		// allows grpc clients to discover the definition of the server without having the protos
		reflection.Register(s)

		err = grpc.RunServer(ctx, s, port)

		if err != nil {
			log.Fatal("error with RunServer", err)
		}
	},
}

func init() {
	rootCommand.AddCommand(serverCommand)

	serverFlags = flags.New("server", serverCommand)

	serverFlags.RegisterInt("port", "p", defaultGRPCPort, "Port of the server", "PORT")

	serverFlags.RegisterString("service-name", "", "hello-grpc", "Name of the application", "APP_NAME")
	serverFlags.RegisterString("environment", "", "local", "Environment of the application", "APP_ENV")
}

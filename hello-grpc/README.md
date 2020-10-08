# hello-grpc

## Introduction

Boilerplate for a simple GRPC application

The associated .proto for this service can be found here:

https://github.com/9count/protos/blob/master/boilerplates/hellogrpc/v1/hello_api.proto

## Getting Started

Run the GRPC server:
```shell script
go run main.go server
```

In a separate terminal, run the GRPC client "Hello" request:
```shell script
go run main.go hello <name>
```

## Testing
```shell script
go test ./...
```

## Building
```shell script
go build
```

## Docker

Build the container locally:
```shell script
GITHUB_DEPLOY_KEY=$(exec cat ~/.ssh/id_rsa) docker build --build-arg GITHUB_DEPLOY_KEY="$GITHUB_DEPLOY_KEY" -t hello-grpc .
```

Run the container:
```shell script
docker run -p "5000:5000" hello-grpc
```

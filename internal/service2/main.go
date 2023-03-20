package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/chiennguyen196/go-opentelemetry-tracing-demo/pkg/genproto/service2"
	"github.com/chiennguyen196/go-opentelemetry-tracing-demo/pkg/tracing"
)

func main() {
	cleanFn := tracing.Init()
	defer func() {
		if err := cleanFn(); err != nil {
			log.Println("Clean tracing error", err)
		}
	}()

	address := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
	RunGRPCServerOnAddr(address, func(server *grpc.Server) {
		svc := GrpcServer{}
		service2.RegisterExampleServiceServer(server, svc)
	})
}

func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	registerServer(server)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Starting: gRPC Listener", addr)
	log.Fatalln(server.Serve(listen))
}

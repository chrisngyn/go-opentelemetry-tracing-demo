package main

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/chiennguyen196/go-opentelemetry-tracing-demo/pkg/genproto/service2"
)

type HttpServer struct {
	client service2.ExampleServiceClient
}

func NewHttpServer() HttpServer {
	conn, err := grpc.Dial("localhost:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatal("Cannot connect to service 2", err)
	}

	client := service2.NewExampleServiceClient(conn)

	return HttpServer{client: client}
}

type Response struct {
	Value string `json:"value"`
}

func (h HttpServer) GetSomething(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.GetSomething(r.Context(), &service2.GetSomethingRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Respond(w, r, Response{Value: resp.Value})

}

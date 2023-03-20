package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/chiennguyen196/go-opentelemetry-tracing-demo/pkg/genproto/service2"
)

type GrpcServer struct {
}

func (g GrpcServer) GetSomething(_ context.Context, _ *service2.GetSomethingRequest) (*service2.GetSomethingResponse, error) {
	if simulateError() {
		return nil, errors.New("something wrong")
	}

	return &service2.GetSomethingResponse{Value: fmt.Sprintf("%d", rand.Int())}, nil
}

func simulateError() bool {
	return rand.Intn(2) == 1
}

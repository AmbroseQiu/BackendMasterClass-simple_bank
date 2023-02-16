package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func GrpcLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Println("receive grpc request")
	result, err := handler(ctx, req)
	return result, err
}

package main

import (
	"context"
	golang_server_side_streaming "github.com/server-side-streaming"
	"github.com/server-side-streaming/notificationservice"
	"github.com/server-side-streaming/notificationservice/notificationproto"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", notificationservice.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	ctx := context.Background()
	redisClient := golang_server_side_streaming.NewRedisClient(ctx)

	handler := notificationservice.NewHandler(redisClient)
	notificationproto.RegisterNotificationServiceServer(grpcServer, handler)

	slog.Info("Listening on " + notificationservice.Address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

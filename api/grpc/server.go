package grpc

import (
	"fmt"
	"kafka-demo/api/grpc/pb"
	db "kafka-demo/internal"
	user "kafka-demo/pkg/user"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
    grpcServer *grpc.Server
)

func Serve(lis net.Listener) error {
	if lis == nil {
		return fmt.Errorf("listener is nil")
	}

	// Create and register a new gRPC server with telemetry
	grpcServer = grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

	// Register your gRPC service
	pb.RegisterUserServiceServer(grpcServer, NewUserServer(user.NewUserService(db.NewUsersRepositoryMongo(db.GetCollection("users")))))

	// Register reflection for gRPC service discovery
	reflection.Register(grpcServer)

	// Start the server
	return grpcServer.Serve(lis)
}


func Stop() {
    grpcServer.GracefulStop()
}
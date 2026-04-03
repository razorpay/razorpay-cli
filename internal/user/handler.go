package user

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userv1 "github.com/razorpay/go-foundation-v2/rpc/go_foundation_v2/user/v1"
)

// GRPCHandler registers the user service with the gRPC server
//
// Parameters:
//   - server: The gRPC server to register the user service with
//
// Returns:
//   - error: An error if the user service registration fails
//   - nil: If the user service registration succeeds
func (s *Server) GRPCHandler(server *grpc.Server) error {
	userv1.RegisterUserServiceServer(server, s)

	return nil
}

// HTTPHandler returns a function that registers the user service HTTP handlers
// with the provided ServeMux for REST API endpoints via gRPC-Gateway.
//
// Parameters:
//   - ctx: Context for request lifecycle management
//
// Returns:
//   - func: A function that takes a ServeMux and address to register handlers
func (s *Server) HTTPHandler(
	ctx context.Context,
) func(mux *runtime.ServeMux, address string) error {
	return func(mux *runtime.ServeMux, address string) error {
		return userv1.RegisterUserServiceHandlerFromEndpoint(
			ctx,
			mux,
			address,
			[]grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			},
		)
	}
}

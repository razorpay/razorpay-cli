package user

import (
	"context"
	"fmt"
	"net/http"

	"buf.build/go/protovalidate"
	"github.com/razorpay/foundation/deps"
	foundationErr "github.com/razorpay/foundation/errors"
	"github.com/razorpay/foundation/grpcserver/mux"

	goutilsError "github.com/razorpay/goutils/errors"

	errpkg "github.com/razorpay/go-foundation-v2/internal/errors"
	"github.com/razorpay/go-foundation-v2/internal/user/model"
	userv1 "github.com/razorpay/go-foundation-v2/rpc/go_foundation_v2/user/v1"
)

// Service defines the interface for user-related business
// logic operations.
type Service interface {
	// CreateUser creates a new user with the provided
	// request data. It returns the created user model
	// or an error if the operation fails.
	//
	// Parameters:
	//   - ctx: Context for request lifecycle management
	//   - req: CreateRequest containing user data to be created
	//
	// Returns:
	//   - *model.User: The created user model if successful
	//   - errors.IError: Error if user creation fails or user already exists
	CreateUser(
		ctx context.Context,
		req *userv1.CreateRequest,
	) (*model.User, goutilsError.IError)

	// GetUser retrieves a user based on the provided
	// request parameters. It returns the user model
	// or an error if the user is not found or operation fails.
	//
	// Parameters:
	//   - ctx: Context for request lifecycle management
	//   - req: GetRequest containing the user ID to retrieve
	//
	// Returns:
	//   - *model.User: The retrieved user model if found, nil if not found
	//   - errors.IError: Error if user retrieval fails or invalid ID provided
	GetUser(
		ctx context.Context,
		req *userv1.GetRequest,
	) (*model.User, goutilsError.IError)
}

// Server provides gRPC server implementation for user
// service operations. It handles request validation,
// business logic delegation, and response formatting
// for all user-related API endpoints.
type Server struct {
	// validator is used to validate the request payload
	validator protovalidate.Validator

	// telemetry is used for logging and metrics collection
	telemetry *deps.Telemetry

	// service is used to delegate the business logic to
	// the service layer
	service Service

	// embedded for forward compatibility
	userv1.UnimplementedUserServiceServer
}

// New creates a new Server instance with the provided dependencies.
// It initializes the protocol buffer validator and sets up the server
// with telemetry and service dependencies.
//
// Parameters:
//   - telemetry: Telemetry instance for logging and metrics
//   - service: Service implementation for business logic operations
//
// Returns:
//   - *Server: Configured server instance ready to handle requests
//   - error: Error if validator initialization fails
func New(telemetry *deps.Telemetry, service Service) (*Server, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("initialize protovalidator: %w", err)
	}

	return &Server{
		validator: validator,
		telemetry: telemetry,
		service:   service,
	}, nil
}

// Create handles user creation requests by validating the input and delegating
// to the service layer for business logic processing.
//
// Parameters:
//   - ctx: Context for request lifecycle management
//   - req: CreateRequest containing user data to create
//
// Returns:
//   - *userv1.User: The created user in protobuf format
//   - error: Error if validation fails or user creation fails
func (s *Server) Create(
	ctx context.Context,
	req *userv1.CreateRequest,
) (*userv1.User, error) {
	if err := s.validator.Validate(req); err != nil {
		ierr := errpkg.ErrValidationFailure.
			New(errpkg.CodeRequestValidationFailure).
			Wrap(err).
			WithPublicMetadata(
				foundationErr.FromProtoValidationError(err),
			)

		return &userv1.User{
			Error: errpkg.ToProtoPublicError(ierr),
		}, errpkg.WrapError(ierr)
	}

	user, err := s.service.CreateUser(ctx, req)
	if err != nil {
		return &userv1.User{
			Error: errpkg.ToProtoPublicError(err),
		}, errpkg.WrapError(err)
	}

	return user.ToProto(), mux.SetCustomStatusCode(ctx, http.StatusCreated)
}

// Get handles user retrieval requests by validating the input and delegating
// to the service layer for business logic processing.
//
// Parameters:
//   - ctx: Context for request lifecycle management
//   - req: GetRequest containing user identifier to retrieve
//
// Returns:
//   - *userv1.User: The retrieved user in protobuf format, or nil if not found
//   - error: Error if validation fails or user retrieval fails
func (s *Server) Get(
	ctx context.Context,
	req *userv1.GetRequest,
) (*userv1.User, error) {
	if err := s.validator.Validate(req); err != nil {
		ierr := errpkg.ErrValidationFailure.
			New(errpkg.CodeRequestValidationFailure).
			Wrap(err).
			WithPublicMetadata(
				foundationErr.FromProtoValidationError(err),
			)

		return &userv1.User{
			Error: errpkg.ToProtoPublicError(ierr),
		}, errpkg.WrapError(ierr)
	}

	user, err := s.service.GetUser(ctx, req)
	if err != nil {
		return &userv1.User{
			Error: errpkg.ToProtoPublicError(err),
		}, errpkg.WrapError(err)
	}

	if user == nil {
		return nil, nil
	}

	return user.ToProto(), nil
}

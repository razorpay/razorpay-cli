package service

import (
	"context"

	"github.com/razorpay/foundation/deps"

	"github.com/razorpay/goutils/errors"
	"github.com/razorpay/goutils/spine"

	"github.com/razorpay/go-foundation-v2/internal/config"
	errpkg "github.com/razorpay/go-foundation-v2/internal/errors"
	"github.com/razorpay/go-foundation-v2/internal/user/model"
	userv1 "github.com/razorpay/go-foundation-v2/rpc/go_foundation_v2/user/v1"
)

// Service provides user-related business logic operations.
type Service struct {
	// repo handles data persistence operations
	repo Repo

	// telemetry provides logging and monitoring capabilities
	telemetry *deps.Telemetry

	awsConfig config.AWSConfig
}

// New creates a new Service instance with the provided telemetry
// and repository dependencies.
//
// Parameters:
//   - telemetry: Telemetry instance for logging and monitoring
//   - repo: Repository instance for data persistence operations
//
// Returns:
//   - *Service: Configured Service instance ready for user operations
func New(telemetry *deps.Telemetry, repo Repo, awsConfig config.AWSConfig) *Service {
	return &Service{
		repo:      repo,
		telemetry: telemetry,
		awsConfig: awsConfig,
	}
}

// CreateUser creates a new user in the system based on the provided request.
//
// Parameters:
//   - ctx: Context for request lifecycle management
//   - req: CreateRequest containing user data to be created
//
// Returns:
//   - *model.User: The created user model if successful
//   - errors.IError: Error if user creation fails or user already exists
func (s *Service) CreateUser(
	ctx context.Context,
	req *userv1.CreateRequest,
) (*model.User, errors.IError) {
	user := model.New()
	user.FromProto(req)

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, spine.UniqueConstraintViolation) {
			ierr := errpkg.ErrBadRequestError.
				New(errpkg.CodeBadRequestError).
				Wrap(err)

			return nil, ierr
		}

		ierr := errpkg.ErrInternalServerError.
			New(errpkg.CodeInternalServerError).
			Wrap(err)

		return nil, ierr
	}

	return user, nil
}

// GetUser retrieves a user from the system based on the provided request.
//
// Parameters:
//   - ctx: Context for request lifecycle management
//   - req: GetRequest containing the user ID to retrieve
//
// Returns:
//   - *model.User: The retrieved user model if found, nil if not found
//   - errors.IError: Error if user retrieval fails or invalid ID provided
func (s *Service) GetUser(
	ctx context.Context,
	req *userv1.GetRequest,
) (*model.User, errors.IError) {
	user := model.New()
	if err := user.FromPublicID(req.GetId()); err != nil {
		ierr := errpkg.ErrInternalServerError.
			New(errpkg.CodeInternalServerError).
			Wrap(err)

		return nil, ierr
	}

	user, err := s.repo.Get(ctx, user.ID)
	if err != nil {
		if errors.Is(err, spine.RecordNotFound) {
			ierr := errpkg.ErrNotFound.
				New(errpkg.CodeNoRecordFound).
				Wrap(err)

			return nil, ierr
		}

		ierr := errpkg.ErrInternalServerError.
			New(errpkg.CodeInternalServerError).
			Wrap(err)

		return nil, ierr
	}

	return user, nil
}

package repo

import (
	"context"

	"github.com/razorpay/foundation/deps"

	"github.com/razorpay/goutils/errors"
	"github.com/razorpay/goutils/spine"
	"github.com/razorpay/goutils/spine/db"

	"github.com/razorpay/go-foundation-v2/internal/user/model"
)

// Repo provides database operations for brand entities.
// It encapsulates database access patterns and includes telemetry
// for observability.
type Repo struct {
	database *spine.Repo

	// Telemetry instance from foundation framework
	// to record observability data (logs, metrics, traces)
	telemetry *deps.Telemetry
}

// New creates a new repository instance with the provided dependencies.
// It initializes the repository with telemetry and database connections.
//
// Parameters:
//   - telemetry: Telemetry instance for observability and metrics collection
//   - database: Database instance for executing queries and transactions
//
// Returns:
//   - *Repo: Configured repository instance ready for database operations
func New(telemetry *deps.Telemetry, database *db.DB) *Repo {
	return &Repo{
		database:  &spine.Repo{Db: database},
		telemetry: telemetry,
	}
}

// Create inserts a new user record into the database.
// It uses the provided context for database operations and returns
// any database errors that occur during the creation process.
//
// Parameters:
//   - ctx: Context for database operation lifecycle and cancellation
//   - user: User model instance to be created in the database
//
// Returns:
//   - errors.IError: Database error if creation fails, nil on success
func (r *Repo) Create(ctx context.Context, user *model.User) errors.IError {
	tx := r.database.DBInstance(ctx).Create(user)

	return spine.GetDBError(tx)
}

// Get retrieves a user record from the database by public ID.
// It converts the public ID to internal ID format and queries the database
// for the corresponding user record.
//
// Parameters:
//   - ctx: Context for database operation lifecycle and cancellation
//   - id: Public ID string of the user to retrieve
//
// Returns:
//   - *model.User: User model instance if found
//   - errors.IError: Database error if query fails or user not found,
//     nil on success
func (r *Repo) Get(
	ctx context.Context,
	id string,
) (*model.User, errors.IError) {
	user := &model.User{}
	if err := user.FromPublicID(id); err != nil {
		return nil, errors.New("invalid public ID")
	}

	tx := r.database.
		DBInstance(ctx).
		Where("id = ?", user.ID).
		First(&user)

	return user, spine.GetDBError(tx)
}

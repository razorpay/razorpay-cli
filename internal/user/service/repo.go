package service

import (
	"context"

	goutilsError "github.com/razorpay/goutils/errors"

	"github.com/razorpay/go-foundation-v2/internal/user/model"
)

// Repo provides database operations for user entities.
// It defines the contract for user data persistence operations
// including creation and retrieval of user records.
type Repo interface {
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
	Create(ctx context.Context, user *model.User) goutilsError.IError

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
	Get(ctx context.Context, id string) (*model.User, goutilsError.IError)
}

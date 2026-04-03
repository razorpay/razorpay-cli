package model

import (
	"strings"

	"github.com/razorpay/goutils/errors"
	"github.com/razorpay/goutils/spine"

	userv1 "github.com/razorpay/go-foundation-v2/rpc/go_foundation_v2/user/v1"
)

const (
	// db table name for users entity
	tableName = "users"
	// entity name for users entity
	entityName = "users"
	// entity prefix used for ID generation and identification
	entityPrefix = "user_"
)

// User represents the user entity
type User struct {
	// embedded model providing common fields like id, created_at,
	// updated_at, and deleted_at.
	spine.SoftDeletableModel

	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// New creates and returns a new User instance with default values.
// This is the recommended way to create new User objects.
//
// Returns:
//   - *User: New User instance ready for use
func New() *User {
	return &User{}
}

// TableName returns the database table name for the [User] model.
// This method is used by GORM to determine which table to use for
// database operations.
//
// Returns:
//   - string: Database table name ("users")
func (*User) TableName() string {
	return tableName
}

// EntityName returns the entity name for the [User] model.
// This is used for entity identification and logging purposes.
//
// Returns:
//   - string: Entity name ("users")
func (*User) EntityName() string {
	return entityName
}

// EntityPrefix returns the entity prefix for the [User] model.
// This is used for ID generation and other entity-related operations.
//
// Returns:
//   - string: Entity prefix ("user_")
func (*User) EntityPrefix() string {
	return entityPrefix
}

// FromPublicID sets the internal ID from a public ID by removing the
// entity prefix. Public IDs are formatted as "user_<internal_id>" and
// this method extracts the internal part.
//
// Note: If the id doesn't have a prefix, it will be used directly
//
// Parameters:
//   - id: Public ID string in the format "user_<internal_id>"
func (u *User) FromPublicID(id string) error {
	if len(id) == 14 {
		u.ID = id
		return nil
	}

	parts := strings.Split(id, entityPrefix)
	if len(parts) != 2 {
		return errors.New("invalid public ID")
	}
	u.ID = parts[1]

	return nil
}

// SetDefaults sets default values for the [User] model.
// Currently no defaults are set, but this method provides a hook for
// future default value logic.
//
// Returns:
//   - errors.IError: Always returns nil as no validation errors can occur
func (u *User) SetDefaults() errors.IError {
	return nil
}

// ToProto converts the internal [User] model to its protocol
// buffer representation.
//
// Returns:
//   - *usersv1.User: Protocol buffer representation of the brand
func (u *User) ToProto() *userv1.User {
	return &userv1.User{
		Id:       u.ID,
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
	}
}

// FromProto populates the internal [User] model from its protocol
// buffer representation.
//
// Parameters:
//   - pb: Protocol buffer representation of the user
func (u *User) FromProto(pb *userv1.CreateRequest) {
	u.Username = pb.GetUsername()
	u.Password = pb.GetPassword()
	u.Email = pb.GetEmail()
}

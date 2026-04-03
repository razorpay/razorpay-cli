package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/razorpay/go-foundation-v2/internal/user/model"
)

func TestNew(t *testing.T) {
	user := model.New()
	assert.NotNil(t, user)
	assert.IsType(t, &model.User{}, user)
}

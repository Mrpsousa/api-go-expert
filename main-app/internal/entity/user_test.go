package entity_test

import (
	"testing"

	"github.com/mrpsousa/api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := entity.NewUser("Roger", "email@email", "12345678")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "Roger", user.Name)
	assert.Equal(t, "email@email", user.Email)
}

func TestUserValidatePassword(t *testing.T) {
	user, err := entity.NewUser("Roger", "email@email", "12345678")
	assert.Nil(t, err)
	assert.NotEqual(t, "12345678", user.Password)
	assert.True(t, user.ValidatePassword("12345678"))
	assert.False(t, user.ValidatePassword("12345"))
}

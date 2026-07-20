package tests

import (
	"strings"
	"testing"

	"villainrsty-ecommerce-server/internal/adapters/security/jwt/service"
	"villainrsty-ecommerce-server/internal/core/shared/models"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	jwtService := service.NewJWTService("test-secret")

	user := &models.User{
		ID:    models.ID("user-123"),
		Email: "userjhondoe@gmail.com",
		Name:  "Jhon Doe",
	}

	token, err := jwtService.GenerateAccessToken(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parts := strings.Split(token, ".")
	assert.Equal(t, len(parts), 3)
}

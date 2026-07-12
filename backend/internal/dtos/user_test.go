package dtos

import (
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewUserResponse(t *testing.T) {
	createdAt := time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 15, 30, 0, 0, time.UTC)
	userID := uuid.Must(uuid.NewV7()).String()

	user := &models.User{
		BaseModel: models.BaseModel{
			ID:        userID,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		Email:           "user@example.com",
		FirstName:       "Jane",
		LastName:        "Doe",
		EmailVerified:   true,
		IsPlatformAdmin: true,
	}

	resp := NewUserResponse(user)
	require.Equal(t, userID, resp.ID)
	require.Equal(t, "user@example.com", resp.Email)
	require.Equal(t, "Jane", resp.FirstName)
	require.Equal(t, "Doe", resp.LastName)
	require.True(t, resp.EmailVerified)
	require.True(t, resp.IsPlatformAdmin)
	require.False(t, resp.MFAEnabled)
	require.Equal(t, createdAt, resp.CreatedAt)
	require.Equal(t, updatedAt, resp.UpdatedAt)
}

package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogoutUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	useCase := NewLogoutUseCase()

	result := useCase.Execute(ctx)

	require.NotNil(t, result)
	assert.Equal(t, "Logged out successfully", result.Message)
}

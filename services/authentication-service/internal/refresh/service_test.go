//go:build unit

package refresh_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
	refreshmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh/mocks"
)

var (
	logger  = zap.NewNop()
	errRepo = errors.New("repository error")
)

func TestService_Generate(t *testing.T) {
	tests := []struct {
		name           string
		input          refresh.GenerateTokenInput
		mocksSetup     func(repo *refreshmocks.MockRepository)
		expectedOutput refresh.GenerateTokenOutput
		expectedErr    error
	}{
		{
			name: "when there is an error storing the refresh token, then it propagates the error",
			input: refresh.GenerateTokenInput{
				UserID: "fake-user-id",
				Role:   "fake-role",
			},
			mocksSetup: func(repo *refreshmocks.MockRepository) {
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(refresh.Token{}, errRepo)
			},
			expectedOutput: refresh.GenerateTokenOutput{},
			expectedErr:    errRepo,
		},
		{
			name: "when the refresh token is generated and stored, then it returns the token",
			input: refresh.GenerateTokenInput{
				UserID: "fake-user-id",
				Role:   "fake-role",
			},
			mocksSetup: func(repo *refreshmocks.MockRepository) {
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, params refresh.CreateTokenParams) (refresh.Token, error) {
						require.Equal(t, "fake-user-id", params.UserID)
						require.Equal(t, "fake-role", params.Role)
						require.NotEmpty(t, params.Token)
						require.NotEmpty(t, params.ExpiresAt)

						// Returning a fake-token to simplify the assertion
						return refresh.Token{Token: "fake-token"}, nil
					})
			},
			expectedOutput: refresh.GenerateTokenOutput{RefreshToken: "fake-token"},
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := refreshmocks.NewMockRepository(ctrl)
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo)
			}

			service := refresh.NewService(logger, repo)
			output, err := service.Generate(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}

//go:build integration

package refresh_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/infraestructure/mongodb"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
)

type refreshRepositoryTestCase[P, W any] struct {
	name            string
	insertDocuments func(t *testing.T, coll *mongo.Collection)
	params          P
	want            W
	wantErr         error
}

func TestRepository_Create(t *testing.T) {
	var (
		now       = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		expiresAt = time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC)
	)
	logger, _ := log.NewTest()

	tests := []refreshRepositoryTestCase[refresh.CreateTokenParams, refresh.Token]{
		{
			name: "when the refresh token already exists, then it should return the error",
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, refresh.Token{
					UserID:    "fake-user-id",
					Role:      "fake-role",
					Token:     "fake-token",
					Status:    refresh.TokenStatusActive,
					ExpiresAt: expiresAt,
					CreatedAt: now,
					UpdatedAt: now,
					DeviceInfo: refresh.DeviceInfo{
						DeviceID:    "fake-device-id",
						UserAgent:   "fake-user-agent",
						IP:          "fake-ip",
						FirstUsedAt: now,
						LastUsedAt:  now,
					},
				})
			},
			params: refresh.CreateTokenParams{
				UserID:    "fake-user-id",
				Role:      "fake-role",
				Token:     "fake-token",
				ExpiresAt: expiresAt,
				Device: refresh.DeviceInfo{
					DeviceID:    "fake-device-id",
					UserAgent:   "fake-user-agent",
					IP:          "fake-ip",
					FirstUsedAt: now,
					LastUsedAt:  now,
				},
			},
			want:    refresh.Token{},
			wantErr: refresh.ErrRefreshTokenAlreadyExists,
		},
		{
			name: "when the refresh token is stored successfully, then it should return the stored token",
			params: refresh.CreateTokenParams{
				UserID:    "fake-user-id",
				Role:      "fake-role",
				Token:     "fake-token",
				ExpiresAt: expiresAt,
				Device: refresh.DeviceInfo{
					DeviceID:    "fake-device-id",
					UserAgent:   "fake-user-agent",
					IP:          "fake-ip",
					FirstUsedAt: now,
					LastUsedAt:  now,
				},
			},
			want: refresh.Token{
				UserID:    "fake-user-id",
				Role:      "fake-role",
				Token:     "fake-token",
				Status:    refresh.TokenStatusActive,
				ExpiresAt: expiresAt,
				CreatedAt: now,
				UpdatedAt: now,
				DeviceInfo: refresh.DeviceInfo{
					DeviceID:    "fake-device-id",
					UserAgent:   "fake-user-agent",
					IP:          "fake-ip",
					FirstUsedAt: now,
					LastUsedAt:  now,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tdb := mongodb.NewTestDB(t)
			defer tdb.Close(t)

			coll := setupTestRefreshTokenCollection(t, tdb.DB)
			if tt.insertDocuments != nil {
				tt.insertDocuments(t, coll)
			}

			repo := refresh.NewRepository(logger, tdb.DB, clock.FixedClock{FixedTime: now})
			token, err := repo.Create(context.Background(), tt.params)

			// Error assertion
			assert.ErrorIs(t, err, tt.wantErr)

			// Validating the token only if there is no error expected
			if tt.wantErr == nil {
				// As the ID is generated by MongoDB, we just check that it is not empty
				assert.NotEmpty(t, token.ID, "ID should not be empty")

				// Doing this as in that way, I can do a direct equal assertion between the token and the expected, so
				// I can detect changes in the struct while ignoring the ID value.
				tt.want.ID = token.ID
				assert.Equal(t, tt.want, token)
			}
		})
	}
}

func TestRepository_Create_UnexpectedFailure(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tdb := mongodb.NewTestDB(t)
	setupTestRefreshTokenCollection(t, tdb.DB)

	repo := refresh.NewRepository(logger, tdb.DB, clock.FixedClock{FixedTime: now})

	// Simulating an unexpected failure by closing the opened connection
	tdb.Close(t)

	_, err := repo.Create(context.Background(), refresh.CreateTokenParams{})
	assert.Error(t, err, "Expected an error due to unexpected failure")
}

func TestRepository_FindActiveToken(t *testing.T) {
	var (
		now       = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		expiresAt = time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC)
		expiredAt = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	)
	logger, _ := log.NewTest()

	tests := []refreshRepositoryTestCase[string, refresh.Token]{
		{
			name: "when the refresh token does not exist, then it should return a refresh token not found error",
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, refresh.Token{
					UserID:    "fake-user-id",
					Role:      "fake-role",
					Token:     "active-token",
					Status:    refresh.TokenStatusActive,
					ExpiresAt: expiresAt,
					CreatedAt: now,
					UpdatedAt: now,
					DeviceInfo: refresh.DeviceInfo{
						DeviceID:    "fake-device-id",
						UserAgent:   "fake-user-agent",
						IP:          "fake-ip",
						FirstUsedAt: now,
						LastUsedAt:  now,
					},
				})
			},
			params:  "unexisting-token",
			want:    refresh.Token{},
			wantErr: refresh.ErrRefreshTokenNotFound,
		},
		{
			name: "when the refresh token is revoked, then it should return a refresh token not found error",
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, refresh.Token{
					UserID:    "fake-user-id",
					Role:      "fake-role",
					Token:     "revoked-token",
					Status:    refresh.TokenStatusRevoked,
					ExpiresAt: expiresAt,
					CreatedAt: now,
					UpdatedAt: now,
					DeviceInfo: refresh.DeviceInfo{
						DeviceID:    "fake-device-id",
						UserAgent:   "fake-user-agent",
						IP:          "fake-ip",
						FirstUsedAt: now,
						LastUsedAt:  now,
					},
				})
			},
			params:  "revoked-token",
			want:    refresh.Token{},
			wantErr: refresh.ErrRefreshTokenNotFound,
		},
		{
			name: "when the refresh token is expired, then it should return a refresh token not found error",
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, refresh.Token{
					UserID:    "fake-user-id",
					Role:      "fake-role",
					Token:     "expired-token",
					Status:    refresh.TokenStatusActive,
					ExpiresAt: expiredAt,
					CreatedAt: now,
					UpdatedAt: now,
					DeviceInfo: refresh.DeviceInfo{
						DeviceID:    "fake-device-id",
						UserAgent:   "fake-user-agent",
						IP:          "fake-ip",
						FirstUsedAt: now,
						LastUsedAt:  now,
					},
				})
			},
			params:  "expired-token",
			want:    refresh.Token{},
			wantErr: refresh.ErrRefreshTokenNotFound,
		},
		{
			name: "when the refresh token is active, then it should return the token",
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, refresh.Token{
					UserID:    "fake-user-id",
					Role:      "fake-role",
					Token:     "active-token",
					Status:    refresh.TokenStatusActive,
					ExpiresAt: expiresAt,
					CreatedAt: now,
					UpdatedAt: now,
					DeviceInfo: refresh.DeviceInfo{
						DeviceID:    "fake-device-id",
						UserAgent:   "fake-user-agent",
						IP:          "fake-ip",
						FirstUsedAt: now,
						LastUsedAt:  now,
					},
				})
			},
			params: "active-token",
			want: refresh.Token{
				UserID:    "fake-user-id",
				Role:      "fake-role",
				Token:     "active-token",
				Status:    refresh.TokenStatusActive,
				ExpiresAt: expiresAt,
				CreatedAt: now,
				UpdatedAt: now,
				DeviceInfo: refresh.DeviceInfo{
					DeviceID:    "fake-device-id",
					UserAgent:   "fake-user-agent",
					IP:          "fake-ip",
					FirstUsedAt: now,
					LastUsedAt:  now,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tdb := mongodb.NewTestDB(t)
			defer tdb.Close(t)

			coll := setupTestRefreshTokenCollection(t, tdb.DB)
			if tt.insertDocuments != nil {
				tt.insertDocuments(t, coll)
			}

			repo := refresh.NewRepository(logger, tdb.DB, clock.FixedClock{FixedTime: now})
			token, err := repo.FindActiveToken(context.Background(), tt.params)

			// Error assertion
			assert.ErrorIs(t, err, tt.wantErr)

			// Validating the token only if there is no error expected
			if tt.wantErr == nil {
				// As the ID is generated by MongoDB, we just check that it is not empty
				assert.NotEmpty(t, token.ID, "ID should not be empty")

				// Doing this as in that way, I can do a direct equal assertion between the token and the expected, so
				// I can detect changes in the struct while ignoring the ID value.
				tt.want.ID = token.ID
				assert.Equal(t, tt.want, token)
			}
		})
	}
}

func TestRepository_FindActiveToken_UnexpectedFailure(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tdb := mongodb.NewTestDB(t)
	setupTestRefreshTokenCollection(t, tdb.DB)

	repo := refresh.NewRepository(logger, tdb.DB, clock.FixedClock{FixedTime: now})

	// Simulating an unexpected failure by closing the opened connection
	tdb.Close(t)

	_, err := repo.FindActiveToken(context.Background(), "")
	assert.Error(t, err, "Expected an error due to unexpected failure")
}

func TestRepository_Expire(t *testing.T) {
	var (
		now          = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		yesterday    = time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
		expiresAt    = time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC)
		newExpiresAt = time.Date(2025, 1, 1, 0, 0, 5, 0, time.UTC)
	)
	logger, _ := log.NewTest()

	tests := []refreshRepositoryTestCase[refresh.ExpireParams, refresh.Token]{
		{
			name: "when the refresh token does not exist, then it should return a refresh token not found error",
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, refresh.Token{
					UserID:    "fake-user-id",
					Role:      "fake-role",
					Token:     "active-token",
					Status:    refresh.TokenStatusActive,
					ExpiresAt: expiresAt,
					CreatedAt: now,
					UpdatedAt: now,
					DeviceInfo: refresh.DeviceInfo{
						DeviceID:    "fake-device-id",
						UserAgent:   "fake-user-agent",
						IP:          "fake-ip",
						FirstUsedAt: now,
						LastUsedAt:  now,
					},
				})
			},
			params:  refresh.ExpireParams{Token: "unexisting-token"},
			want:    refresh.Token{},
			wantErr: refresh.ErrRefreshTokenNotFound,
		},
		{
			name: "when the refresh token is revoked, then it should return a refresh token not found error",
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, refresh.Token{
					UserID:    "fake-user-id",
					Role:      "fake-role",
					Token:     "revoked-token",
					Status:    refresh.TokenStatusRevoked,
					ExpiresAt: expiresAt,
					CreatedAt: now,
					UpdatedAt: now,
					DeviceInfo: refresh.DeviceInfo{
						DeviceID:    "fake-device-id",
						UserAgent:   "fake-user-agent",
						IP:          "fake-ip",
						FirstUsedAt: now,
						LastUsedAt:  now,
					},
				})
			},
			params:  refresh.ExpireParams{Token: "revoked-token"},
			want:    refresh.Token{},
			wantErr: refresh.ErrRefreshTokenNotFound,
		},
		{
			name: "when the refresh token is already expired, then it should return a refresh token not found error",
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, refresh.Token{
					UserID:    "fake-user-id",
					Role:      "fake-role",
					Token:     "expired-token",
					Status:    refresh.TokenStatusActive,
					ExpiresAt: newExpiresAt,
					CreatedAt: now,
					UpdatedAt: now,
					DeviceInfo: refresh.DeviceInfo{
						DeviceID:    "fake-device-id",
						UserAgent:   "fake-user-agent",
						IP:          "fake-ip",
						FirstUsedAt: now,
						LastUsedAt:  now,
					},
				})
			},
			params: refresh.ExpireParams{
				Token:     "expired-token",
				ExpiresAt: newExpiresAt,
			},
			want:    refresh.Token{},
			wantErr: refresh.ErrRefreshTokenNotFound,
		},
		{
			name: "when the refresh token is active, then it should return the token expired",
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, refresh.Token{
					UserID:    "fake-user-id",
					Role:      "fake-role",
					Token:     "active-token",
					Status:    refresh.TokenStatusActive,
					ExpiresAt: expiresAt,
					CreatedAt: yesterday,
					UpdatedAt: yesterday,
					DeviceInfo: refresh.DeviceInfo{
						DeviceID:    "fake-device-id",
						UserAgent:   "fake-user-agent",
						IP:          "fake-ip",
						FirstUsedAt: yesterday,
						LastUsedAt:  yesterday,
					},
				})
			},
			params: refresh.ExpireParams{
				Token:     "active-token",
				ExpiresAt: newExpiresAt,
			},
			want: refresh.Token{
				UserID:    "fake-user-id",
				Role:      "fake-role",
				Token:     "active-token",
				Status:    refresh.TokenStatusActive,
				ExpiresAt: newExpiresAt,
				CreatedAt: yesterday,
				UpdatedAt: now,
				DeviceInfo: refresh.DeviceInfo{
					DeviceID:    "fake-device-id",
					UserAgent:   "fake-user-agent",
					IP:          "fake-ip",
					FirstUsedAt: yesterday,
					LastUsedAt:  yesterday,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tdb := mongodb.NewTestDB(t)
			defer tdb.Close(t)

			coll := setupTestRefreshTokenCollection(t, tdb.DB)
			if tt.insertDocuments != nil {
				tt.insertDocuments(t, coll)
			}

			repo := refresh.NewRepository(logger, tdb.DB, clock.FixedClock{FixedTime: now})
			token, err := repo.Expire(context.Background(), tt.params)

			// Error assertion
			assert.ErrorIs(t, err, tt.wantErr)

			// Validating the token only if there is no error expected
			if tt.wantErr == nil {
				// As the ID is generated by MongoDB, we just check that it is not empty
				assert.NotEmpty(t, token.ID, "ID should not be empty")

				// Doing this as in that way, I can do a direct equal assertion between the token and the expected, so
				// I can detect changes in the struct while ignoring the ID value.
				tt.want.ID = token.ID
				assert.Equal(t, tt.want, token)
			}
		})
	}
}

func TestRepository_Expire_UnexpectedFailure(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tdb := mongodb.NewTestDB(t)
	setupTestRefreshTokenCollection(t, tdb.DB)

	repo := refresh.NewRepository(logger, tdb.DB, clock.FixedClock{FixedTime: now})

	// Simulating an unexpected failure by closing the opened connection
	tdb.Close(t)

	_, err := repo.Expire(context.Background(), refresh.ExpireParams{})
	assert.Error(t, err, "Expected an error due to unexpected failure")
}

func setupTestRefreshTokenCollection(t *testing.T, db *mongo.Database) *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.Collection(refresh.CollectionName)

	// Create unique index on token.
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: refresh.FieldToken, Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := coll.Indexes().CreateOne(ctx, indexModel); err != nil {
		t.Fatalf("Failed to create unique index: %v", err)
	}

	return coll
}

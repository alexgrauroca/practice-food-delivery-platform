package refresh

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/logctx"
)

// TokenStatus represents the state or condition of a token, typically used to indicate its validity or current status.
type TokenStatus string

const (
	// CollectionName defines the name of the database collection used to store refresh tokens.
	CollectionName = "refresh_tokens"

	// TokenStatusActive represents a token that is currently valid and in use.
	TokenStatusActive TokenStatus = "active"
	// TokenStatusRevoked represents a token that has been invalidated and can no longer be used.
	TokenStatusRevoked TokenStatus = "revoked"

	// FieldToken represents the database field name for storing token values.
	FieldToken = "token"
	// FieldUserID defines the database field name for storing user identifier values.
	FieldUserID = "user_id"
)

// Repository defines a contract for storing and managing refresh tokens in a persistence layer.
//
//go:generate mockgen -destination=./mocks/repository_mock.go -package=refresh_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh Repository
type Repository interface {
	Create(ctx context.Context, params CreateTokenParams) (Token, error)
}

// Token represents a token used to refresh authentication credentials for a specific user and role.
type Token struct {
	ID        string      `bson:"_id,omitempty"`
	UserID    string      `bson:"user_id"`
	Role      string      `bson:"role"`
	Token     string      `bson:"token"`
	Status    TokenStatus `bson:"status"`
	ExpiresAt time.Time   `bson:"expires_at"`
	CreatedAt time.Time   `bson:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at"`
}

// CreateTokenParams defines the parameters required to create a new token for a user.
type CreateTokenParams struct {
	UserID    string
	Role      string
	Token     string
	ExpiresAt time.Time
}

type repository struct {
	logger     *zap.Logger
	collection *mongo.Collection
	clock      clock.Clock
}

// NewRepository creates a new Repository instance.
func NewRepository(logger *zap.Logger, db *mongo.Database, clk clock.Clock) Repository {
	return &repository{
		logger:     logger,
		collection: db.Collection(CollectionName),
		clock:      clk,
	}
}

func (r *repository) Create(ctx context.Context, params CreateTokenParams) (Token, error) {
	now := r.clock.Now()
	token := Token{
		UserID:    params.UserID,
		Role:      params.Role,
		Token:     params.Token,
		ExpiresAt: params.ExpiresAt,
		CreatedAt: now,
		UpdatedAt: now,
		Status:    TokenStatusActive,
	}

	res, err := r.collection.InsertOne(ctx, token)
	if err != nil {
		logctx.LoggerWithRequestInfo(ctx, r.logger).Error("Failed to store refresh token", zap.Error(err))
		return Token{}, err
	}

	token.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return token, nil
}

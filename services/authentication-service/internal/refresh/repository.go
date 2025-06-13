package refresh

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/logctx"
)

const (
	// CollectionName defines the name of the database collection used to store refresh tokens.
	CollectionName = "refresh_tokens"

	// FieldToken represents the database field name for storing token values.
	FieldToken = "token"
	// FieldStatus represents the database field name for storing token status information.
	FieldStatus = "status"
	// FieldExpiresAt represents the database field name for storing the expiration time of a token.
	FieldExpiresAt = "expires_at"
)

// Repository defines a contract for storing and managing refresh tokens in a persistence layer.
//
//go:generate mockgen -destination=./mocks/repository_mock.go -package=refresh_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh Repository
type Repository interface {
	Create(ctx context.Context, params CreateTokenParams) (Token, error)
	FindActiveToken(ctx context.Context, refreshToken string) (Token, error)
}

// CreateTokenParams defines the parameters required to create a new token for a user.
type CreateTokenParams struct {
	UserID    string
	Role      string
	Token     string
	Device    DeviceInfo
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
		UserID:     params.UserID,
		Role:       params.Role,
		Token:      params.Token,
		Status:     TokenStatusActive,
		DeviceInfo: params.Device,
		ExpiresAt:  params.ExpiresAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	token.DeviceInfo.FirstUsedAt = now
	token.DeviceInfo.LastUsedAt = now

	res, err := r.collection.InsertOne(ctx, token)
	if err != nil {
		logctx.LoggerWithRequestInfo(ctx, r.logger).Error("Failed to store refresh token", zap.Error(err))
		return Token{}, err
	}

	token.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return token, nil
}

func (r *repository) FindActiveToken(ctx context.Context, refreshToken string) (Token, error) {
	token := Token{}
	searchParams := bson.M{
		FieldToken:  refreshToken,
		FieldStatus: TokenStatusActive,
		FieldExpiresAt: bson.M{
			"$gt": r.clock.Now(),
		},
	}
	err := r.collection.FindOne(ctx, searchParams).Decode(&token)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logctx.LoggerWithRequestInfo(ctx, r.logger).Warn("Refresh token not found")
			return Token{}, ErrRefreshTokenNotFound
		}
		logctx.LoggerWithRequestInfo(ctx, r.logger).Error("Failed to find active refresh refreshToken", zap.Error(err))
		return Token{}, err
	}

	return token, nil
}

package refresh

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/infraestructure/mongodb"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
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
	// FieldUpdatedAt represents the database field name for storing the timestamp of the last update.
	FieldUpdatedAt = "updated_at"
)

// Repository defines a contract for storing and managing refresh tokens in a persistence layer.
//
//go:generate mockgen -destination=./mocks/repository_mock.go -package=refresh_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh Repository
type Repository interface {
	Create(ctx context.Context, params CreateTokenParams) (Token, error)
	FindActiveToken(ctx context.Context, refreshToken string) (Token, error)
	Expire(ctx context.Context, params ExpireParams) (Token, error)
}

type repository struct {
	logger     log.Logger
	collection *mongo.Collection
	clock      clock.Clock
}

// NewRepository creates a new Repository instance.
func NewRepository(logger log.Logger, db *mongo.Database, clk clock.Clock) Repository {
	return &repository{
		logger:     logger,
		collection: db.Collection(CollectionName),
		clock:      clk,
	}
}

// CreateTokenParams defines the parameters required to create a new token for a user.
type CreateTokenParams struct {
	UserID    string
	Role      string
	TenantID  string
	Token     string
	Device    DeviceInfo
	ExpiresAt time.Time
}

func (r *repository) Create(ctx context.Context, params CreateTokenParams) (Token, error) {
	logger := r.logger.WithContext(ctx)

	now := r.clock.Now()
	token := Token{
		UserID:     params.UserID,
		Role:       params.Role,
		TenantID:   params.TenantID,
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
		if mongodb.IsDuplicateKeyError(err) {
			logger.Error("Duplicate refresh token", err)
			return Token{}, ErrRefreshTokenAlreadyExists
		}
		logger.Error("Failed to store refresh token", err)
		return Token{}, err
	}

	token.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return token, nil
}

func (r *repository) FindActiveToken(ctx context.Context, refreshToken string) (Token, error) {
	logger := r.logger.WithContext(ctx)

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
			logger.Warn("Refresh token not found")
			return Token{}, ErrRefreshTokenNotFound
		}
		logger.Error("Failed to find active refresh refreshToken", err)
		return Token{}, err
	}

	return token, nil
}

// ExpireParams defines parameters needed to update a token's expiration status.
// Token represents the refresh token to be expired.
// ExpiresAt specifies the new expiration time for the token.
type ExpireParams struct {
	Token     string
	ExpiresAt time.Time
}

func (r *repository) Expire(ctx context.Context, params ExpireParams) (Token, error) {
	logger := r.logger.WithContext(ctx)

	var token Token
	filter := bson.M{
		FieldToken:  params.Token,
		FieldStatus: TokenStatusActive,
		// If the token was already expired, then we do nothing
		FieldExpiresAt: bson.M{
			"$gt": params.ExpiresAt,
		},
	}

	update := bson.M{
		"$set": bson.M{
			FieldExpiresAt: params.ExpiresAt,
			FieldUpdatedAt: r.clock.Now(),
		},
	}

	// Returning the updated document
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&token)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Warn("Refresh token not found")
			return Token{}, ErrRefreshTokenNotFound
		}
		logger.Error("Failed to expire refresh token", err)
		return Token{}, err
	}
	return token, nil
}

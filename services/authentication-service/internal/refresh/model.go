package refresh

import "time"

// TokenStatus represents the state or condition of a token, typically used to indicate its validity or current status.
type TokenStatus string

const (
	// TokenStatusActive represents a token that is currently valid and in use.
	TokenStatusActive TokenStatus = "active"
	// TokenStatusRevoked represents a token that has been invalidated and can no longer be used.
	TokenStatusRevoked TokenStatus = "revoked"
)

// DeviceInfo represents information about a device.
type DeviceInfo struct {
	DeviceID    string    `bson:"device_id"`
	UserAgent   string    `bson:"user_agent"`
	IP          string    `bson:"ip"`
	FirstUsedAt time.Time `bson:"first_used_at"`
	LastUsedAt  time.Time `bson:"last_used_at"`
}

// Token represents a token used to refresh authentication credentials for a specific user and role.
type Token struct {
	ID         string      `bson:"_id,omitempty"`
	UserID     string      `bson:"user_id"`
	Role       string      `bson:"role"`
	Token      string      `bson:"token"`
	Status     TokenStatus `bson:"status"`
	DeviceInfo DeviceInfo  `bson:"device_info"`
	ExpiresAt  time.Time   `bson:"expires_at"`
	CreatedAt  time.Time   `bson:"created_at"`
	UpdatedAt  time.Time   `bson:"updated_at"`
}

package domain

import (
	"time"
)

type TokenRevoke struct {
	ID        string    `bson:"_id,omitempty"`
	Token     string    `bson:"token"`
	UserID    string    `bson:"user_id"`
	RevokedAt time.Time `bson:"revoked_at"`
}

type LoginAttempt struct {
	Email     string    `bson:"email"`
	Attempts  int       `bson:"attempts"`
	LastTry   time.Time `bson:"last_try"`
	BlockedAt time.Time `bson:"blocked_at,omitempty"`
}

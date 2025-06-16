package mongo

import (
	"context"
	"time"

	"github.com/nightnice1st/testGridWhiz/internal/auth/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthRepository struct {
	db          *mongo.Database
	tokenColl   *mongo.Collection
	attemptColl *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) *AuthRepository {
	return &AuthRepository{
		db:          db,
		tokenColl:   db.Collection("tokenRevoke"),
		attemptColl: db.Collection("loginAttempts"),
	}
}

func (r *AuthRepository) RevokeToken(token, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	revoke := &domain.TokenRevoke{
		Token:     token,
		UserID:    userID,
		RevokedAt: time.Now(),
	}

	_, err := r.tokenColl.InsertOne(ctx, revoke)
	return err
}

func (r *AuthRepository) IsTokenRevoked(token string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result domain.TokenRevoke
	err := r.tokenColl.FindOne(ctx, bson.M{"token": token}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *AuthRepository) RecordLoginAttempt(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": email}
	update := bson.M{
		"$inc": bson.M{"attempts": 1},
		"$set": bson.M{"last_try": time.Now()},
	}

	_, err := r.attemptColl.UpdateOne(ctx, filter, update,
		&options.UpdateOptions{Upsert: &[]bool{true}[0]})

	return err
}

func (r *AuthRepository) GetLoginAttempts(email string) (*domain.LoginAttempt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var attempt domain.LoginAttempt
	err := r.attemptColl.FindOne(ctx, bson.M{"email": email}).Decode(&attempt)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &attempt, err
}

func (r *AuthRepository) ResetLoginAttempts(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.attemptColl.DeleteOne(ctx, bson.M{"email": email})
	return err
}

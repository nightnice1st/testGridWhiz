package mongo

import (
	"context"
	"time"

	"github.com/nightnice1st/testGridWhiz/internal/users/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
	coll := db.Collection("users")

	// Create indexes for better performance
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
	}

	coll.Indexes().CreateMany(ctx, indexes)

	return &userRepository{
		db:   db,
		coll: coll,
	}
}

func (r *userRepository) Create(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid.Hex()
	}

	return nil
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = r.coll.FindOne(ctx, bson.M{"_id": oid, "deleted_at": nil}).Decode(&user)
	if err != nil {
		return nil, err
	}

	user.ID = oid.Hex()
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User
	err := r.coll.FindOne(ctx, bson.M{"email": email, "deleted_at": nil}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}

	user.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name": user.Name,
			// "email":      user.Email,
			"updated_at": user.UpdatedAt,
		},
	}

	_, err = r.coll.UpdateOne(ctx, bson.M{"_id": oid}, update)
	return err
}

func (r *userRepository) SoftDelete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}

	_, err = r.coll.UpdateOne(ctx, bson.M{"_id": oid}, update)
	return err
}

func (r *userRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.coll.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *userRepository) List(page, limit int, nameFilter, emailFilter string) ([]*domain.User, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"deleted_at": nil}
	if nameFilter != "" {
		filter["name"] = bson.M{"$regex": nameFilter, "$options": "i"}
	}
	if emailFilter != "" {
		filter["email"] = bson.M{"$regex": emailFilter, "$options": "i"}
	}

	// Count total documents
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip
	skip := (page - 1) * limit

	// Find documents with pagination
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			continue
		}

		if oid, ok := cursor.Current.Lookup("_id").ObjectIDOK(); ok {
			user.ID = oid.Hex()
		}

		users = append(users, &user)
	}

	return users, int(total), nil
}

package auth

import (
	"context"
	"errors"

	"github.com/google/uuid" 
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	app_auth "github.com/AldiandyaIrsyad/author-notes/internal/auth" 
)

// mongoAuthRepository implements the AuthRepository interface using MongoDB.
type mongoAuthRepository struct {
	collection *mongo.Collection
}

// NewMongoAuthRepository creates a new instance of mongoAuthRepository.
func NewMongoAuthRepository(db *mongo.Database) app_auth.AuthRepository {
	return &mongoAuthRepository{
		collection: db.Collection("users"), // Assuming the collection name is "users"
	}
}

// CreateUser inserts a new user into the database.
func (r *mongoAuthRepository) CreateUser(ctx context.Context, user *app_auth.User) error {
	// Generate a new UUID if the ID is empty
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	_, err := r.collection.InsertOne(ctx, user)
	// Consider handling potential duplicate key errors specifically
	// if a race condition occurs despite the pre-check in the service.
	return err
}

// FindByUsername retrieves a user by their username.
func (r *mongoAuthRepository) FindByUsername(ctx context.Context, username string) (*app_auth.User, error) {
	var user app_auth.User
	filter := bson.M{"username": username}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, app_auth.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmailOrUsername retrieves a user by their email or username.
func (r *mongoAuthRepository) FindByEmailOrUsername(ctx context.Context, email, username string) (*app_auth.User, error) {
	var user app_auth.User
	filter := bson.M{
		"$or": []bson.M{
			{"email": email},
			{"username": username},
		},
	}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, app_auth.ErrUserNotFound // Use the domain error
		}
		return nil, err
	}
	return &user, nil
}

// FindByID retrieves a user by their ID.
func (r *mongoAuthRepository) FindByID(ctx context.Context, id string) (*app_auth.User, error) {
	var user app_auth.User
	// Query using the string ID directly, assuming it's stored as a string (UUID)
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, app_auth.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

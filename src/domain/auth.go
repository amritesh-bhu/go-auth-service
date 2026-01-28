package domain

import (
	"context"
	"errors"
	"time"

	"github.com/go-auth-service/src/entities"
	"github.com/go-auth-service/src/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginRequest struct {
	EmailId  string `bson:"emailId" json:"emailId"`
	Password string `bson:"password" json:"password"`
}

type RegisterResponse struct {
	ID        primitive.ObjectID `json:"id"`
	EmailId   string             `json:"emailId"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
}

var UsersCollection *mongo.Collection

func RegisterUser(ctx context.Context, user *entities.User) (*RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	hashedPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	result, err := UsersCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	// assign generated _id back to user
	response := &RegisterResponse{
		ID:        result.InsertedID.(primitive.ObjectID),
		EmailId:   user.EmailId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return response, nil
}

func LoginUser(ctx context.Context, cred *LoginRequest) (*entities.User, error) {
	var user entities.User
	if err := UsersCollection.FindOne(ctx, bson.M{"emailId": cred.EmailId}).Decode(&user); err != nil {
		return nil, errors.New("User does not exist!")
	}

	if helper.ComparePassword(user.Password, cred.Password) {
		return nil, errors.New("Invalid credentials!")
	}

	return &user, nil
}

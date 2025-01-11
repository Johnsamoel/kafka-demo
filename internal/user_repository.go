package db

import (
	"context"
	"fmt"
	model "kafka-demo/pkg/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersRepositoryMongo struct {
	collection *mongo.Collection
}

func NewUsersRepositoryMongo(collection *mongo.Collection) *UsersRepositoryMongo {
	return &UsersRepositoryMongo{collection: collection}
}

func (u *UsersRepositoryMongo) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {

	_, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		fmt.Println("create user Error: %+v", err)
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return user, nil
}

func (u *UsersRepositoryMongo) UpdateUser(ctx context.Context, user *model.User, userId string) error {
	update := bson.M{
		"$set": bson.M{
			"first_Name": user.FirstName,
			"last_Name":  user.LastName,
			"email":      user.Email,
			"isNotified": user.IsNotified,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := u.collection.FindOneAndUpdate(ctx, bson.M{"userId": userId}, update, opts)

	err := result.Decode(user)
	if err != nil {
		fmt.Printf("update user Error: %+v\n", err)
		if err == mongo.ErrNoDocuments {
			return status.Errorf(codes.NotFound, "user not found")
		}
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}

func (u *UsersRepositoryMongo) GetUserById(ctx context.Context, userId string) (*model.User, error) {

	user := &model.User{}

	err := u.collection.FindOne(ctx, bson.M{"userId": userId}).Decode(user)
	if err != nil {
		fmt.Println("find user Error: %+v", err)
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return user, nil
}


func (u *UsersRepositoryMongo) GetUserByEmail(ctx context.Context, userEmail string) (*model.User, error) {

	user := &model.User{}

	err := u.collection.FindOne(ctx, bson.M{"email": userEmail}).Decode(user)
	if err != nil {
		fmt.Println("find user Error: %+v", err)
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return user, nil
}


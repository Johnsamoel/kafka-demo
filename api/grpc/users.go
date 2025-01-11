package grpc

import (
	"context"
	"encoding/json"
	"kafka-demo/api/grpc/pb"
	"kafka-demo/pkg/kafka"
	model "kafka-demo/pkg/user"
	"kafka-demo/utils"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExportedUserServer struct {
	pb.UnimplementedUserServiceServer
	userService model.UserService
}

func NewUserServer(userService model.UserService) *ExportedUserServer {
	return &ExportedUserServer{
		userService: userService,
	}
}

func (u *ExportedUserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	newUser := model.User{
		UserId:     "",
		FirstName:  req.User.FirstName,
		LastName:   req.User.LastName,
		Email:      req.User.Email,
		Password:   req.User.Password,
		IsNotified: false,
	}

	// Validate User struct
	validate := validator.New()
	err := validate.Struct(newUser)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// hash user password
	hashedPassword, err := utils.HashPassword(req.User.Password)

	if err != nil {
		return nil, status.Error(codes.Internal, "something went wrong!")
	}

	MappedUser := mapCreateUserProtoToModel(req.User)

	MappedUser.Password = hashedPassword
	MappedUser.UserId = uuid.NewString()

	userData, err := u.userService.CreateUser(ctx, MappedUser)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// create kafka msg
	UserRegisterationMsg := &kafka.KafkaNotificationMsg{
		User:               userData,
		NotificationStatus: kafka.NEW_USER_WELCOMING_NOTIFICATIONS,
	}

	convertedMsg, err := json.Marshal(UserRegisterationMsg)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	HeadersMap := make(map[string]string)
	HeadersMap["retries"] = "3"

	err = kafka.ProduceMessage(kafka.Notifications, userData.UserId, string(convertedMsg), HeadersMap)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateUserResponse{
		User: mapUserModelToProto(userData),
	}, nil
}

func (u *ExportedUserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	err := uuid.Validate(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	userData, err := u.userService.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetUserResponse{User: mapUserModelToProto(userData)}, nil
}

func (u *ExportedUserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	newUser := model.User{
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
		Email:     req.User.Email,
		Password:  req.User.Password,
	}

	validate := validator.New()

	err := validate.Struct(newUser)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID are required")
	}

	userData := mapUserProtoToModel(req.User)

	err = u.userService.UpdateUser(ctx, userData, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateUserResponse{User: mapUserModelToProto(userData)}, nil
}

func mapUserProtoToModel(userData *pb.User) *model.User {
	return &model.User{
		FirstName:  userData.FirstName,
		LastName:   userData.LastName,
		Email:      userData.Email,
		Password:   userData.Password,
		IsNotified: userData.IsNotified,
	}
}

func mapCreateUserProtoToModel(userData *pb.OmmittedUser) *model.User {
	return &model.User{
		FirstName:  userData.FirstName,
		LastName:   userData.LastName,
		Email:      userData.Email,
		Password:   userData.Password,
		IsNotified: false,
	}
}

func mapUserModelToProto(userData *model.User) *pb.User {
	return &pb.User{
		UserId:     userData.UserId,
		FirstName:  userData.FirstName,
		LastName:   userData.LastName,
		Email:      userData.Email,
		Password:   userData.Password,
		IsNotified: userData.IsNotified,
	}
}

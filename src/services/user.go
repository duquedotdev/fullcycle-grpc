package services

import (
	"context"
	"fmt"
	"github.com/duquedotdev/fullcycle-grpc/pb"
	"io"
	"log"
	"time"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {

	fmt.Println(req.Name)

	return &pb.User{
		Id:    "1",
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil

}

func (*UserService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {

	err := stream.Send(&pb.UserResultStream{
		Status: "Starting to create a new user...",
		User:   &pb.User{},
	})
	if err != nil {
		return err
	}

	time.Sleep(3 * time.Second)

	err = stream.Send(&pb.UserResultStream{
		Status: "User created with success!",
		User: &pb.User{
			Id:    "1",
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})
	if err != nil {
		return err
	}

	return nil

}

func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {

	users := []*pb.User{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}
		if err != nil {
			log.Fatalf("error when receiving the stream: %v", err)
		}

		users = append(users, &pb.User{
			Id:    req.GetId(),
			Name:  req.GetName(),
			Email: req.GetEmail(),
		})

		fmt.Println("User added: ", req.GetName())

	}
}
func (*UserService) AddUsersBothStream(stream pb.UserService_AddUsersBothStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error when receiving the stream: %v", err)
		}
		err = stream.Send(&pb.UserResultStream{
			Status: "User added with success!",
			User:   req,
		})
		if err != nil {
			log.Fatalf("error when sending the stream: %v", err)
		}
	}
}

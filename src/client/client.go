package main

import (
	"context"
	"fmt"
	"github.com/duquedotdev/fullcycle-grpc/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {

	connection, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	AddUsersBothStream(client)

}

func AddUser(client pb.UserServiceClient) {

	request := &pb.User{
		Id:    "1",
		Name:  "Duque",
		Email: "felipe@duque.dev",
	}

	response, err := client.AddUser(context.Background(), request)
	if err != nil {
		log.Fatalf("error when calling AddUser: %v", err)
	}

	fmt.Println("User added: ", response)

}

func AddUserVerbose(client pb.UserServiceClient) {

	request := &pb.User{
		Id:    "1",
		Name:  "Duque",
		Email: "felipe@duque.dev",
	}

	response, err := client.AddUserVerbose(context.Background(), request)
	if err != nil {
		log.Fatalf("error when calling AddUser: %v", err)
	}

	for {
		stream, err := response.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error when receiving the stream: %v", err)
		}
		fmt.Println("Status: ", stream.Status)
		fmt.Println("User: ", stream.User)
	}

}

func AddUsers(client pb.UserServiceClient) {

	request := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "Duque",
			Email: "felipe@duque.dev",
		},
		&pb.User{
			Id:    "2",
			Name:  "Duque",
			Email: "felipe@duque.dev",
		},
		&pb.User{
			Id:    "3",
			Name:  "Duque",
			Email: "felipe@duque.dev",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("error when calling AddUsers: %v", err)
	}

	for _, user := range request {
		fmt.Println("Sending user: ", user)
		err := stream.Send(user)
		if err != nil {
			return
		}
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error when receiving the response: %v", err)
	}

	fmt.Println("Users added: ", res)
}

func AddUsersBothStream(client pb.UserServiceClient) {

	stream, err := client.AddUsersBothStream(context.Background())

	if err != nil {
		log.Fatalf("error when calling AddUsers: %v", err)
	}

	request := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "Duque",
			Email: "felipe@duque.dev",
		},
		&pb.User{
			Id:    "2",
			Name:  "Duque",
			Email: "felipe@duque.dev",
		},
		&pb.User{
			Id:    "3",
			Name:  "Duque",
			Email: "felipe@duque.dev",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range request {
			fmt.Println("Sending user: ", req)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error when sending user: %v", err)
			}
			time.Sleep(1000 * time.Millisecond)

		}
		err = stream.CloseSend()
		if err != nil {
			return
		}
	}()

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error when receiving the stream: %v", err)
			}
			fmt.Println("Status: ", response.Status)
			fmt.Println("User: ", response.User)
		}
		close(wait)
	}()

	<-wait

}

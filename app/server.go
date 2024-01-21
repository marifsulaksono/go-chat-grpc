package main

import (
	"go-chat-grpc/protobuf/chatserver"
	"go-chat-grpc/services"
	"log"
	"net"

	"google.golang.org/grpc"
)

const port = ":50505"

func main() {
	// init listener
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error listening port : %v", err)
	}

	grpcServer := grpc.NewServer() // gRPC server instance
	cs := services.ChatServer{}
	chatserver.RegisterChatServicesServer(grpcServer, &cs)

	log.Println("grpc server started at", listener.Addr())
	err = grpcServer.Serve(listener) // gRPC listen and serve
	if err != nil {
		log.Fatalf("Error start grpc server : %v", err)
	}
}

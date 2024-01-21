package main

import (
	"bufio"
	"context"
	"fmt"
	"go-chat-grpc/protobuf/chatserver"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Enter server IP:Port ::: ")
	reader := bufio.NewReader(os.Stdin)
	serverId, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read server port from console : %v", err)
	}

	serverId = strings.Trim(serverId, "\r\n")
	log.Println("Connecting : " + serverId)

	conn, err := grpc.Dial(serverId, grpc.WithInsecure()) // connect to gRPC server
	if err != nil {
		log.Fatalf("Failed to connect to grpc server : %v", err)
	}
	defer conn.Close()

	client := chatserver.NewChatServicesClient(conn) // call ChatService to create a stream
	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to call Chat Service : %v", err)
	}

	// implement communication with gRPC server
	ch := clientHandler{stream: stream}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	// blocker
	bl := make(chan bool)
	<-bl
}

// client handle
type clientHandler struct {
	stream     chatserver.ChatServices_ChatServiceClient
	clientName string
}

func (ch *clientHandler) clientConfig() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Your name : ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read from console : %v", err)
	}
	ch.clientName = strings.Trim(name, "\r\n")
}

// send message
func (ch *clientHandler) sendMessage() {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read from console : %v", err)
		}

		message = strings.Trim(message, "\r\n")
		clientMessageBox := &chatserver.FromClient{
			Name: ch.clientName,
			Body: message,
		}

		err = ch.stream.Send(clientMessageBox)
		if err != nil {
			log.Printf("Error while sending message to server : %v", err)
		}
	}
}

// receive message
func (ch *clientHandler) receiveMessage() {
	for {
		msg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error when receive message : %v", err)
		}

		fmt.Printf("%s : %s \n", msg.Name, msg.Body)
	}
}

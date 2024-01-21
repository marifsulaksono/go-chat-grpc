package services

import (
	"log"
	"math/rand"
	"sync"
	"time"

	pbChatServ "go-chat-grpc/protobuf/chatserver"
)

// define a messageUnit struct for handing messages in server
type messageUnit struct {
	ClientName        string
	MessageBody       string
	MessageUniqueCode int
	ClientUniqueCode  int
}

// define messageHandle struct to hold slice of messageUnits
type messageHandle struct {
	MessageQue []messageUnit
	mut        sync.Mutex
}

type ChatServer struct{}

var messageHandleObj = messageHandle{} // for handling asynchronous read write operation

func (c *ChatServer) ChatService(cs pbChatServ.ChatServices_ChatServiceServer) error {
	clientUniqueCode := rand.Intn(1e6)
	errch := make(chan error)

	// init a go routine to asynchronous run receiveFromStream methods
	go receiveFromStream(cs, clientUniqueCode, errch)

	// init a go routine to asynchronous run receiveFromStream methods
	go sendToStream(cs, clientUniqueCode, errch)

	return <-errch
}

func receiveFromStream(cs pbChatServ.ChatServices_ChatServiceServer, clientUniqueCode int, errch chan error) {
	// create infinite loop for continously  receive message from client
	for {
		msg, err := cs.Recv()
		if err != nil {
			log.Printf("Error when receiving message from client : %v", err)
			errch <- err
		} else {
			messageHandleObj.mut.Lock()
			messageHandleObj.MessageQue = append(messageHandleObj.MessageQue, messageUnit{
				ClientName:        msg.Name,
				MessageBody:       msg.Body,
				MessageUniqueCode: rand.Intn(1e8),
				ClientUniqueCode:  clientUniqueCode,
			})

			messageHandleObj.mut.Unlock()

			log.Printf("%v", messageHandleObj.MessageQue[len(messageHandleObj.MessageQue)-1])
		}
	}
}

func sendToStream(cs pbChatServ.ChatServices_ChatServiceServer, clientUniqueCode int, errch chan error) {
	for {
		for {
			time.Sleep(500 * time.Millisecond)
			messageHandleObj.mut.Lock()

			if len(messageHandleObj.MessageQue) == 0 {
				messageHandleObj.mut.Unlock()
				break
			}

			senderNameForClient := messageHandleObj.MessageQue[0].ClientName
			senderUniqueCode := messageHandleObj.MessageQue[0].ClientUniqueCode
			messageForClient := messageHandleObj.MessageQue[0].MessageBody

			messageHandleObj.mut.Unlock()

			// send message to designated client (do not send to the same client)
			if senderUniqueCode != clientUniqueCode {
				err := cs.Send(&pbChatServ.FromServer{
					Name: senderNameForClient,
					Body: messageForClient,
				})

				if err != nil {
					errch <- err
				}

				messageHandleObj.mut.Lock()

				if len(messageHandleObj.MessageQue) > 1 {
					messageHandleObj.MessageQue = messageHandleObj.MessageQue[1:] // delete the message at index 0 after sending to receiver
				} else {
					messageHandleObj.MessageQue = []messageUnit{}
				}

				messageHandleObj.mut.Unlock()
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

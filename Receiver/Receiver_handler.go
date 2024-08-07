package Receiver

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	network = "tcp"
	address = "127.0.0.1"
	port    = "9080"
)

type fileInfo struct {
	filename   string
	filesize   string
	senderName string
}
type receiver struct {
	listener           net.Listener
	listenAddr         string
	transferConnection net.Conn
}
type Receiver interface {
	startServer()
	acceptConnection()
	authTransfer() (bool, *fileInfo, error)
	receiveFile(file fileInfo, wg *sync.WaitGroup)
}

// Server init
func (r *receiver) startServer() {
	Listener, err := net.Listen(network, address+":"+port)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
	fmt.Println("Server listening on port: " + port + ", [CTRL+C to stop server]")
	r.listener = Listener
}

// Connection handler
func (r *receiver) acceptConnection() {
	conn, err := r.listener.Accept()
	if err != nil {
		log.Fatal("Cannot accept connection: ", err)
	}
	r.transferConnection = conn
}

// File information extraction from sender message
func extractFileInfo(msg string) *fileInfo {
	tmp := strings.Split(msg, ":")
	return &fileInfo{
		filename:   tmp[0],
		filesize:   tmp[1],
		senderName: tmp[2],
	}
}

// Auth handler
func (r *receiver) authTransfer() (bool, *fileInfo, error) {
	// Get file info from sender
	senderMsgBuffer := make([]byte, 120)
	senderMsg, err := r.transferConnection.Read(senderMsgBuffer)
	if err != nil {
		return false, nil, err
	}
	senderMsgString := string(senderMsgBuffer[:senderMsg])
	incomingFileInfo := extractFileInfo(senderMsgString)

	// Gather response
	fmt.Println("Do you want to receive file: " + incomingFileInfo.filename + " ,size: " + incomingFileInfo.filesize + ", from: " + incomingFileInfo.senderName + " (Y/N): ")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Response: ")
	scanner.Scan()
	receiverResponse := strings.ToUpper(scanner.Text())

	// Return to response to sender
	_, err = r.transferConnection.Write([]byte(receiverResponse))
	if err != nil {
		return false, nil, err
	}

	// Save response for later use
	switch receiverResponse {
	case "Y":
		return true, incomingFileInfo, nil
	default:
		return false, nil, nil
	}
}

// File receiving logic
func (r *receiver) receiveFile(file fileInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	defer r.transferConnection.Close()
	outfile, err := os.Create(file.filename)
	if err != nil {
		log.Fatal("Cannot create file: ", err)
	}
	defer outfile.Close()
	incomingFileBuffer := make([]byte, 1024)
	for {
		fileChunks, err := r.transferConnection.Read(incomingFileBuffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Cannot read connection: ", err)
		}
		_, err = outfile.Write(incomingFileBuffer[:fileChunks])
		if err != nil {
			log.Fatal("Cannot write file")
		}
	}
}
func StartReceiver() {
	var wg sync.WaitGroup
	receiver := receiver{
		listenAddr: address + ":" + port,
	}
	receiver.startServer()
	for {
		wg.Add(1)
		receiver.acceptConnection()
		fmt.Println("Accepted connection from: " + receiver.listener.Addr().String())
		receiveFile, incomingFile, err := receiver.authTransfer()
		if err != nil {
			log.Println("Error in authentication :(", err)
		}
		switch receiveFile {
		case true:
			receiver.receiveFile(*incomingFile, &wg)
			fmt.Println("File: " + incomingFile.filename + " received successfully :)")
			wg.Wait()
		case false:
			fmt.Println("Transfer denied")
			wg.Done()
		}

	}
}

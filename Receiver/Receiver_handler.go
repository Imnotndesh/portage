package Receiver

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	port    = "9080"
	network = "tcp"
	address = "127.0.0.1"
)

func ExtractDataFromMessage(message string) (string, string, string) {
	tmpMsg := strings.Split(message, ":")
	return tmpMsg[0], tmpMsg[1], tmpMsg[2]
}
func receiveFile(conn net.Conn, filename string) {
	defer conn.Close()
	outfile, err := os.Create(filename)
	if err != nil {
		log.Panic("Error creating file", err)
	}
	defer outfile.Close()
	buffer := make([]byte, 1024)
	for {
		fileChunks, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Cannot read from connection")
		}
		outfile.Write(buffer[:fileChunks])
	}
	fmt.Println("Received file successfully")
}
func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	var response []byte
	var responseStr string
	buffer := make([]byte, 1024)
	senderMsg, err := conn.Read(buffer)
	if err != nil {
		log.Println("Error reading from connection")
		return
	}
	requestMsg := string(buffer[:senderMsg])
	newFilename, newFileSize, senderAlias := ExtractDataFromMessage(requestMsg)
	fmt.Println("Do you want to receive file : " + newFilename + " , size: " + newFileSize + ", From: " + senderAlias + " : (Y/N)")
	fmt.Println("Response")
	fmt.Scanln(&responseStr)
	if err != nil {
		log.Panic("Cannot read response from user")
	}
	response = []byte(strings.ToUpper(responseStr))
	_, err = conn.Write(response)
	if err != nil {
		log.Panic("Error writing to connection")
		return
	}
	switch responseStr {
	case "y":
		receiveFile(conn, newFilename)
		conn.Close()
		wg.Done()
	default:
		wg.Done()
	}

}
func StartReceiveServer() {
	var conn net.Conn
	var wg sync.WaitGroup
	server, err := net.Listen(network, address+":"+port)
	if err != nil {
		log.Panic("Cannot start server", err)
		return
	}
	defer server.Close()
	fmt.Println("Receiver listening on port: " + port + ", CTRL+C to stop server")
	for {
		wg.Add(1)
		fmt.Println("Waiting for sender...")
		conn, err = server.Accept()
		if err != nil {
			log.Panic("Error accepting connection", err)
			continue
		}
		fmt.Println("----------File transfer----------")
		go handleConnection(conn, &wg)
		wg.Wait()
		fmt.Println("--------------Done*--------------")
		fmt.Println("")
	}
}

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
	machinename = "testreceiver"
	port        = ":9080"
	network     = "tcp"
)

func ExtractDataFromMessage(message string) (string, string) {
	tmpMsg := strings.Split(message, ":")
	return tmpMsg[0], tmpMsg[1]
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
			log.Panic("Cannot read from connection")
		}
		outfile.Write(buffer[:fileChunks])
	}
	fmt.Println("Received file successfully")
}
func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	var response []byte
	var responseStr string
	buffer := make([]byte, 1024)
	senderMsg, err := conn.Read(buffer)
	if err != nil {
		log.Println("Error reading from connection")
		return
	}
	requestMsg := string(buffer[:senderMsg])
	newFilename, senderAlias := ExtractDataFromMessage(requestMsg)
	fmt.Println("Do you want to receive file : " + newFilename + ", From: " + senderAlias + " : (Y/N)")
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
	if responseStr == "N" {
		wg.Done()
	}
	receiveFile(conn, newFilename)
	wg.Done()
}
func StartReceiveServer() {
	server, err := net.Listen(network, port)
	if err != nil {
		log.Panic("Cannot start server", err)
		return
	}
	defer server.Close()
	fmt.Println("Receiver listening on 9080, CTRL+C to stop server")
	var conn net.Conn
	var wg sync.WaitGroup
	for {
		wg.Add(1)
		fmt.Println("Waiting for connection...")
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

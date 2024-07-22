package Sender

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const (
	machinename = "testsender"
	network     = "tcp"
)

func sendFile(conn net.Conn, filename string) {
	defer conn.Close()
	fileBuffer := make([]byte, 1024)
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	for {
		chunks, err := file.Read(fileBuffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write(fileBuffer[:chunks])
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	fmt.Println("Sent successfully to remote")
}
func responseListener(conn net.Conn) (string, net.Conn) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Println(err.Error())
	}
	return string(buffer[:n]), conn
}

var requestMessage []byte

func SendToRemoteMachine(filename string, remoteIP string) {
	receiverIP := remoteIP + ":9080"
	conn, err := net.Dial(network, receiverIP)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	reqString := filename + ":" + machinename
	requestMessage = []byte(reqString)
	_, err = conn.Write(requestMessage)
	if err != nil {
		log.Fatal(err.Error())
	}
	responseText, establishedConn := responseListener(conn)
	switch responseText {
	case "Y":
		sendFile(establishedConn, filename)
		os.Exit(0)
	case "N":
		fmt.Println("Transfer request denied")
		os.Exit(0)
	default:
		fmt.Println("Invalid response: Switching to default (Stop transfer)")
		os.Exit(0)
	}
}

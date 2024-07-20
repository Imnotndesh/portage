package Sender

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
)

func SendFileToIP(filepath string, receiverIP string) error {
	var err error
	file, err := os.Open(filepath)
	if err != nil {
		return errors.New("Error opening file: " + err.Error())
	}
	defer file.Close()
	connection, err := net.Dial("tcp", receiverIP)
	if err != nil {
		return errors.New("Error establishing a connection ")
	}
	defer connection.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error reading file: " + err.Error())
			os.Exit(1)
		}
		_, err = connection.Write(buffer[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

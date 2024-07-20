package Receiver

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func StartReceiverServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	outfile, err := os.Create("./output")
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		outfile.Write(buffer[:n])
	}
	fmt.Println("File received successfully")
}

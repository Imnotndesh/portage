package main

import (
	"fmt"
	"os"
	"portage/Receiver"
	"portage/Sender"
)

const defaultPort = "8756"

func main() {
	switch os.Args[1] {
	case "-s":
		var filepath = os.Args[2]
		var receiverIP = os.Args[3]
		err := Sender.SendFileToIP(filepath, receiverIP)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("File set successfully")
	case "-r":
		var customPort = os.Args[2]
		if customPort == "" {
			fmt.Println("Starting server at port: " + defaultPort)
			Receiver.StartReceiverServer(defaultPort)
		} else {
			fmt.Println("Starting server at port: " + customPort)
			Receiver.StartReceiverServer(customPort)
		}
	case "-h":
		fmt.Println("Help Text placeholder")
	}
}

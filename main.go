package main

import (
	"fmt"
	"os"
	"portage/Receiver"
	"portage/Sender"
	"strings"
)

func main() {
	var modes = strings.ToUpper(os.Args[1])
	switch modes {
	case "-S":
		var receiverIP = os.Args[2]
		files := os.Args[3:]
		Sender.SendFile(receiverIP, files...)
	case "-R":
		Receiver.StartReceiveServer()
	default:
		fallthrough
	case "-H":
		fmt.Println("Help Text placeholder")
	}
}

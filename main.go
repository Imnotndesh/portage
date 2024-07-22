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
		var file = os.Args[2]
		var receiverIP = os.Args[3]
		Sender.SendToRemoteMachine(file, receiverIP)
	case "-R":
		Receiver.StartReceiveServer()
	default:
		fallthrough
	case "-H":
		fmt.Println("Help Text placeholder")
	}
}

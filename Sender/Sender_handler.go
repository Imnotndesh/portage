package Sender

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

const (
	machinename = "testsender"
	port        = "9080"
	network     = "tcp"
)

type fileinfo struct {
	filename     string
	filesize     int64
	fileContents *os.File
}
type sender struct {
	conection   net.Conn
	machinename string
}

func getFileInfo(file string) *fileinfo {
	filestats, err := os.Stat(file)
	if err != nil {
		log.Fatal("cannot open file", err)
	}
	fileconts, err := os.Open(file)
	if err != nil {
		log.Fatal("Cannot open file: ", err)
	}
	return &fileinfo{
		filename:     filestats.Name(),
		filesize:     filestats.Size(),
		fileContents: fileconts,
	}
}
func newSenderConnection(ip string) *sender {
	conn, err := net.Dial(network, ip+":"+port)
	if err != nil {
		log.Fatal("cannot open connection", err)
	}
	return &sender{
		conection:   conn,
		machinename: machinename,
	}
}
func (s *sender) authenticateTransfer(file string) (bool, error) {
	newFileInfo := getFileInfo(file)
	filesize := strconv.FormatInt(newFileInfo.filesize, 10)
	requestString := newFileInfo.filename + ":" + filesize + ":" + s.machinename
	fmt.Println(newFileInfo.filesize)
	requestBuffer := []byte(requestString)
	_, err := s.conection.Write(requestBuffer)
	if err != nil {
		return false, err
	}
	response := func(conn net.Conn) string {
		responseBuffer := make([]byte, 1024)
		response, err := conn.Read(responseBuffer)
		if err != nil {
			log.Fatal("cannot read response", err)
		}
		s.conection = conn
		return string(responseBuffer[:response])
	}(s.conection)
	fmt.Println(response)
	switch response {
	case "Y":
		return true, nil
	default:
		return false, nil
	}

}
func (s *sender) sendToRemote(file string, wg *sync.WaitGroup) {
	fileInfo := getFileInfo(file)
	filebuffer := make([]byte, 1024)
	defer fileInfo.fileContents.Close()
	defer wg.Done()
	for {
		chunks, err := fileInfo.fileContents.Read(filebuffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		_, err = s.conection.Write(filebuffer[:chunks])
		if err != nil {
			log.Fatal("Transfer error:", err)
		}
	}
}
func SendFile(ip string, files ...string) {
	var wg sync.WaitGroup
	for _, file := range files {
		newsender := newSenderConnection(ip)
		respValue, err := newsender.authenticateTransfer(file)
		if err != nil {
			log.Println("Error in authentication", err)
			continue
		}
		switch respValue {
		case true:
			wg.Add(1)
			go newsender.sendToRemote(file, &wg)
			wg.Wait()
			log.Println("File sent to remote: ", ip)
		case false:
			log.Println("Transfer request denied by receiver: ", ip)
		}
		newsender.conection.Close()
	}
}

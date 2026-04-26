package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(file io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer file.Close()
		defer close(out)

		str := ""
		for {
			data := make([]byte, 8)
			n, err := file.Read(data)
			if err != nil {
				break
			}

			data = data[:n]

			if index := bytes.IndexByte(data, '\n'); index != -1 {
				str += string(data[:index])
				data = data[index+1:]
				out <- str
				str = ""
			}

			str += string(data)
		}

		if str != "" {
			out <- str
		}
	}()

	return out
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		fmt.Println("connection has been accepted")
		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Printf("Read: %s\n", line)
		}

	}
}

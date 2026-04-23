package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main () {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("Something went wrong")
	}

	defer file.Close()

	str := ""

	for {
		data := make([]byte, 8)
		n, err := file.Read(data)
		if err != nil {
			break
		}

		data = data[:n]

		if idx := bytes.IndexByte(data, '\n'); idx != -1 {
			str += string(data[:idx])
			data=data[idx+1:]
			fmt.Printf("read:%s\n", str)
			str = ""
		}

		str += string(data)
		
	}
}
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("error on opening listener:", err)
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error on connection:", err)
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("%s\n", line)
		}

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		currentLine := ""
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					ch <- currentLine
					currentLine = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error reading file: %s\n", err)
				break
			}

			str := string(buffer[:n])
			slc := strings.Split(str, "\n")

			for i := range len(slc) - 1 {
				ch <- currentLine + slc[i]
				currentLine = ""
			}

			currentLine += slc[len(slc)-1]
		}
		close(ch)
	}()

	return ch
}

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("message.txt")
	if err != nil {
		log.Fatal("error loading file:", err)
	}

	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
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

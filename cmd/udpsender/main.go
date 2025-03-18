package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("error on creating udp addr: %s\n", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("error on creating udp connection: %s\n", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error on reading stdin: %s", err)
		}
		_, err = conn.Write([]byte(str))
		if err != nil {
			fmt.Printf("error on sending message: %s", err)
		}

	}

}

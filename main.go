package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/juancwu/jellyfand/command"
)

const (
	SERVER_UNIX_SOCKET_PATH = "/tmp/jellyfand.server.sock"
	BUFFER_BLOCK_SIZE       = 4096
)

func main() {
	fmt.Println("Jellyfand")

	socket, err := net.Listen("unix", SERVER_UNIX_SOCKET_PATH)
	if err != nil {
		log.Fatal(err)
	}

	// cleanup channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(SERVER_UNIX_SOCKET_PATH)
		os.Exit(1)
	}()

	// accept connections
	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			defer conn.Close()

			buf := make([]byte, BUFFER_BLOCK_SIZE)

			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			// handle request
			res, err := command.ParseInput(buf[:n])
			_, err = conn.Write(res)
			if err != nil {
				log.Fatal(err)
			}
		}(conn)
	}
}

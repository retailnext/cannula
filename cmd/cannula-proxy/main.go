package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var listenAddr = flag.String("addr", "127.0.0.1:0", "TCP address/port to listen on")

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("Usage: cannula-proxy /path/to/socket")
	}

	sockPath := args[0]

	ln, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	unixConn, err := net.Dial("unix", sockPath)
	if err != nil {
		log.Fatal(err)
	}
	unixConn.Close()

	fmt.Printf("Listening on %s, proxying to %s\n", ln.Addr(), sockPath)

	for {
		tcpConn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(tcpConn net.Conn) {
			defer tcpConn.Close()

			unixConn, err := net.Dial("unix", sockPath)
			if err != nil {
				log.Fatal(err)
			}
			defer unixConn.Close()

			go func() {
				io.Copy(tcpConn, unixConn)
			}()

			io.Copy(unixConn, tcpConn)
		}(tcpConn)
	}
}

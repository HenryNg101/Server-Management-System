package main

import (
	"fmt"
	"log"
	"net"
	"runtime"
)

const (
	startPort  = 10000
	numServers = 10000
)

func main() {
	// Use all CPUs efficiently
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(runtime.NumCPU())

	for i := 0; i < numServers; i++ {
		port := startPort + i

		go func(p int) {
			addr := fmt.Sprintf(":%d", p)
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				log.Printf("Failed to listen on port %d: %v\n", p, err)
				return
			}
			log.Printf("Listening on %d\n", p)

			for {
				conn, err := ln.Accept()
				if err != nil {
					continue
				}

				// Handle connection quickly
				go func(c net.Conn) {
					defer c.Close()
					c.Write([]byte("OK"))
				}(conn)
			}
		}(port)
	}

	select {} // block forever
}

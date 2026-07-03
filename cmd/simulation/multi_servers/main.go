package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime"

	"github.com/HenryNg101/server-management-system/internal/app"
	"github.com/HenryNg101/server-management-system/internal/feature/user"
	"github.com/HenryNg101/server-management-system/internal/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	startPort  = 10000
	numServers = 10000
)

func SeedAdmin(userService user.Service) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	user := model.User{
		Email:    "admin@example.com",
		Password: string(hash),
		Role:     "admin",
	}
	userService.CreateUser(context.Background(), user)
}

func main() {
	// Use all CPUs efficiently
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println(runtime.NumCPU())

	// Create app
	newApplication, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	SeedAdmin(newApplication.UserService)

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

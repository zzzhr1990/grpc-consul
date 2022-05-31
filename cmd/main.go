package main

import (
	"log"
	"time"

	"github.com/zzzhr1990/grpc-consul/server"
)

func main() {
	err := server.RegisterToConsul(&server.ConsulRegisterConfig{
		ConsulAddress: "",
	})
	if err != nil {
		log.Fatalf("cannot register to consul: %v", err)
	}
	for range time.Tick(time.Second) {
		log.Println("tick")
	}
}

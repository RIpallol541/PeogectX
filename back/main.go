package main

import (
	"fmt"

	"github.com/RIpallol541/PeogectX/db"
	"github.com/RIpallol541/PeogectX/udp"
)

func main() {
	// Initialize MongoDB connection
	db.Connect()
	defer db.Disconnect()

	go udp.StartUDPServer()

	fmt.Println("UDP server is running")
	select {} // Keep the program running
}
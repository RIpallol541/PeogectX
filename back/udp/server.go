package udp

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	
	"github.com/RIpallol541/PeogectX/handlers"
)

const UDP_PORT = ":8081"

var sessionStore = struct {
	sync.Mutex
	sessions map[string]bool
}{sessions: make(map[string]bool)}

// StartUDPServer starts the UDP server for handling JSON requests
func StartUDPServer() {
	addr, err := net.ResolveUDPAddr("udp", UDP_PORT)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("Listening on UDP %s\n", UDP_PORT)

	buffer := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP: %v\n", err)
			continue
		}

		go handleUDPRequest(conn, remoteAddr, buffer[:n])
	}
}

func handleUDPRequest(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	var request map[string]interface{}
	if err := json.Unmarshal(data, &request); err != nil {
		fmt.Printf("Invalid JSON: %v\n", err)
		return
	}

	action, ok := request["action"].(string)
	if !ok {
		fmt.Printf("Missing action in request\n")
		return
	}

	switch action {
	case "auth":
		var user struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.Unmarshal(data, &user); err != nil {
			fmt.Printf("Invalid user data: %v\n", err)
			return
		}
		response := handlers.HandleAuth(user)
		if response["status"] == "Authentication successful" {
			sessionStore.Lock()
			sessionStore.sessions[addr.String()] = true
			sessionStore.Unlock()
		}
		sendUDPResponse(conn, addr, response)
	case "matrix_task":
		sessionStore.Lock()
		isAuthorized := sessionStore.sessions[addr.String()]
		sessionStore.Unlock()

		if !isAuthorized {
			response := map[string]string{"error": "Unauthorized"}
			sendUDPResponse(conn, addr, response)
			return
		}

		params := map[string]uint32{
			"sizeMatrix":   uint32(request["sizeMatrix"].(float64)),
			"maxDimension": uint32(request["maxDimension"].(float64)),
		}
		go handlers.PerformMatrixMultiplicationTask(conn, addr, params)
	default:
		response := map[string]string{"error": "Unknown action"}
		sendUDPResponse(conn, addr, response)
	}
}

func sendUDPResponse(conn *net.UDPConn, addr *net.UDPAddr, response interface{}) {
	responseData, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Error marshaling response: %v\n", err)
		return
	}

	_, err = conn.WriteToUDP(responseData, addr)
	if err != nil {
		fmt.Printf("Error sending response: %v\n", err)
	}
}

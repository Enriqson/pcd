package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const CONN_TYPE = "gorpc"

func createUDPConnection(serverAddress string) net.Conn {
	remoteAddr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		os.Exit(1)
	}

	return conn
}

func createTCPConnection(serverAddress string) net.Conn {
	// Dial to the server address
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("failed to connect to server: ", err)
		os.Exit(1)
	}
	return conn
}

func saveResults(rttList []float64) {
	// Save the round-trip times to a JSON file
	filename := "../rtt_times_" + CONN_TYPE + ".json"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(rttList)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		os.Exit(1)
	}

	fmt.Println("Round-trip times saved to " + filename)
}

func makeRequestUDP(conn net.Conn, query_url string) {
	// Send the query to the server
	_, err := conn.Write([]byte(query_url))
	if err != nil {
		fmt.Println("Error sending query:", err)
		os.Exit(1)
	}

	// Set a read deadline to avoid blocking forever
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// Read the response from the server
	buffer := make([]byte, 10*1024) // Buffer to store server response
	n, _, err := conn.(*net.UDPConn).ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error reading response:", err)
		os.Exit(1)
	}
	response := string(buffer[:n])
	_ = response

	// Print the server's response (list of strings)
	//fmt.Println("Received response:")
	//fmt.Println(string(buffer[:n]))
}

func makeRequestTCP(conn net.Conn, query_url string) {
	// Send the query to the server
	_, err := conn.Write([]byte(query_url))
	if err != nil {
		fmt.Println("Error sending query:", err)
		os.Exit(1)
	}

	// Set a read deadline to avoid blocking forever
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	buffer := make([]byte, 10*1024)

	n, err := conn.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading from connection:", err)
		return
	}

	response := string(buffer[:n])
	_ = response
	//fmt.Println(response)
}

func main() {
	// Server address
	serverAddress := "localhost:9090" // Ensure this matches the server address
	var makeRequest func(conn net.Conn, query_url string)
	var conn net.Conn
	// Create a UDP connection
	if CONN_TYPE == "udp" {
		conn = createUDPConnection(serverAddress)
		makeRequest = makeRequestUDP
	} else {
		conn = createTCPConnection(serverAddress)
		makeRequest = makeRequestTCP
	}

	defer conn.Close()

	var rttList []float64
	// The query string to send to the server
	query_url := "https://stackoverflow.com"
	for i := 1; i <= 100; i++ {

		fmt.Println("Sending query:", i, "...")
		startTime := time.Now()

		makeRequest(conn, query_url)

		duration := time.Since(startTime).Seconds()

		rttList = append(rttList, duration)

		//avoid rate limit errors to the crawled site
		time.Sleep(1 * time.Second)
	}
	saveResults(rttList)
}

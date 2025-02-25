package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"io"
	"net"
	"os"

	"github.com/PuerkitoBio/goquery"
)

const CONN_TYPE = "tcp"

func isValidHTTPUrl(urlStr string) bool {
	// Ensure the URL starts with "http://" or "https://"
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return false
	}

	// Parse the URL
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check if the URL's scheme is "http" or "https"
	return parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https"
}

func crawlPage(link string) []string {
	res, err := http.Get(link)
	if err != nil {
		log.Printf("Error fetching URL %s: %v", link, err)
		return []string{}
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("Error parsing document for URL %s: %v", link, err)
		return []string{}
	}

	var links []string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		linkRef := s.AttrOr("href", "")
		if isValidHTTPUrl(linkRef) {
			links = append(links, linkRef)
		}
	})
	if len(links) == 0 {
		println("Error crawling page")
	}
	return links
}

func createUDPConnection(serverAddress string) *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	// Create a UDP connection
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		os.Exit(1)
	}

	return conn
}

func startUDPServer(serverAddress string) {
	conn := createUDPConnection(serverAddress)
	fmt.Println("UDP Server listening on", serverAddress)
	// Buffer to store incoming requests
	buffer := make([]byte, 10*1024)

	for {
		// Read incoming data
		handleUDPRequest(conn, buffer)
	}

}

func handleUDPRequest(conn *net.UDPConn, buffer []byte) {
	n, clientAddr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error reading UDP packet:", err)
		return
	}

	// Extract the query_url from the received message
	query_url := string(buffer[:n])
	fmt.Println("Received query_url:", query_url)

	// Call the crawl function
	results := crawlPage(query_url)

	// Prepare the response (a list of strings)
	response := strings.Join(results, "\n")

	// Send the response back to the client
	_, err = conn.WriteToUDP([]byte(response), clientAddr)
	if err != nil {
		fmt.Println("Error sending response:", err)
	} else {
		fmt.Println("Sent response to client")
	}
}

func createTCPListener(serverAddress string) net.Listener {
	ln, err := net.Listen("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		os.Exit(1)
	}

	return ln
}

func startTCPServer(serverAddress string) {
	ln := createTCPListener(serverAddress)
	fmt.Println("TCP Server listening on", serverAddress)
	buffer := make([]byte, 1024)
	for {
		// Accept a connection
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting TCP connection:", err)
			continue
		}

		// Handle the connection
		go handleTCPConnection(conn, buffer)
	}
}

func handleTCPConnection(conn net.Conn, buffer []byte) {
	defer conn.Close()

	for {
		n, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading from TCP connection:", err)
			return
		}

		if n == 0 {
			// No data to read, break the loop
			break
		}
		query_url := string(buffer[:n])
		// Call the crawl function
		results := crawlPage(query_url)

		// Prepare the response (a list of strings)
		response := strings.Join(results, "\n")
		// Print the received message
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to TCP connection:", err)
			return
		}
		fmt.Println("Sent response to client")
		fmt.Println(response)
	}

}

func main() {
	// Define server address
	serverAddress := "localhost:9090"

	if CONN_TYPE == "udp" {
		startUDPServer(serverAddress)
	} else {
		startTCPServer(serverAddress)
	}

}

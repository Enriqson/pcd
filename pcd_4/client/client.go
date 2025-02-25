package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/valyala/gorpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "pcd/proto"
)

const CONN_TYPE = "grpc"
const CALLS_PER_CLIENT = 50
const NUM_CLIENTS = 1
const SERVER_ADDRESS = "localhost:9090"
const URL_TO_CRAWL = "https://example.com"

func saveResults(rttList []float64) {
	// Save the round-trip times to a JSON file
	filename := "../rtt_times_" + CONN_TYPE + "_" + strconv.Itoa(NUM_CLIENTS) + "_clients" + ".json"
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

func goRpcClient(client_number int, wg *sync.WaitGroup, mx *sync.Mutex, rttListGlobal *[]float64) {
	defer wg.Done()

	// Create an RPC client
	client := &gorpc.Client{
		// TCP address of the server.
		Addr: SERVER_ADDRESS,
	}
	client.Start()
	defer client.Stop()

	var rttList []float64
	for i := 1; i <= CALLS_PER_CLIENT; i++ {
		fmt.Println("Sending query:", i, " in client ", client_number, "...")

		startTime := time.Now()
		response, err := client.Call(URL_TO_CRAWL)
		_ = response
		if err != nil {
			log.Fatalf("Error calling CrawlURL: %v", err)
		}

		duration := time.Since(startTime).Seconds()

		rttList = append(rttList, duration)

		//avoid rate limit errors to the crawled site
		//time.Sleep(100 * time.Millisecond)
	}

	mx.Lock()
	*rttListGlobal = append(*rttListGlobal, rttList...)
	mx.Unlock()
}

func grpcClient(client_number int, wg *sync.WaitGroup, mx *sync.Mutex, rttListGlobal *[]float64) {
	defer wg.Done()

	// Connect to the gRPC server
	conn, err := grpc.NewClient(SERVER_ADDRESS, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCrawlerServiceClient(conn)

	// Call the CrawlURL method
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rttList []float64
	for i := 1; i <= CALLS_PER_CLIENT; i++ {
		fmt.Println("Sending query:", i, " in client ", client_number, "...")

		startTime := time.Now()

		response, err := client.CrawlURL(ctx, &pb.CrawlRequest{Url: URL_TO_CRAWL})
		_ = response
		if err != nil {
			log.Fatalf("Error calling CrawlURL: %v", err)
		}

		duration := time.Since(startTime).Seconds()

		rttList = append(rttList, duration)

		//avoid rate limit errors to the crawled site
		//time.Sleep(100 * time.Millisecond)
	}

	mx.Lock()
	*rttListGlobal = append(*rttListGlobal, rttList...)
	mx.Unlock()
}

func main() {
	wg := sync.WaitGroup{}
	mx := sync.Mutex{}
	var rttListGlobal []float64
	var clientFunc func(client_number int, wg *sync.WaitGroup, mx *sync.Mutex, rttListGlobal *[]float64)
	// The query string to send to the server

	if CONN_TYPE == "grpc" {
		clientFunc = grpcClient
	} else {
		clientFunc = goRpcClient
	}

	for i := 1; i <= NUM_CLIENTS; i++ {
		wg.Add(1)
		go clientFunc(i, &wg, &mx, &rttListGlobal)
	}
	wg.Wait()
	fmt.Println(rttListGlobal)
	saveResults(rttListGlobal)
}

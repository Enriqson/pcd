package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"

	"net"

	"github.com/PuerkitoBio/goquery"

	"github.com/valyala/gorpc"
	"google.golang.org/grpc"

	// Replace with the path to your generated proto files
	pb "pcd/proto"
)

const CONN_TYPE = "grpc"

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

// Server struct
type Server struct {
	pb.UnimplementedCrawlerServiceServer
}

// CrawlURL implementation
func (s *Server) CrawlURL(ctx context.Context, req *pb.CrawlRequest) (*pb.CrawlResponse, error) {
	//log.Printf("Crawling URL in GRPC: %s", req.GetUrl())
	results := crawlPage(req.Url)

	return &pb.CrawlResponse{Urls: results}, nil
}

func grpcServer() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCrawlerServiceServer(grpcServer, &Server{})

	log.Println("GRPC Server is running...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func goRpcServer() {

	s := &gorpc.Server{
		// Accept clients on this TCP address.
		Addr: ":9090",

		// Echo handler - just return back the message we received from the client
		Handler: func(clientAddr string, request interface{}) interface{} {
			request_url := request.(string)
			log.Printf("Crawling URL in GoRPC: %s", request_url)
			results := crawlPage(request_url)
			return results
		},
	}
	log.Println("GoRPC Server is running...")
	if err := s.Serve(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}
}

func main() {
	var serverFunc func()

	if CONN_TYPE == "grpc" {
		serverFunc = grpcServer
	} else {
		serverFunc = goRpcServer
	}

	serverFunc()
}

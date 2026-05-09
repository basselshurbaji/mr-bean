package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	serverURL := os.Getenv("MR_BEAN_SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080"
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN environment variable is required")
	}

	client := NewClient(serverURL, token)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mr-bean",
		Version: "1.0.0",
	}, nil)

	registerBeanTools(server, client)
	registerGearTools(server, client)
	registerExtractionTools(server, client)
	registerUserTools(server, client)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

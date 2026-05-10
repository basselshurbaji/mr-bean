package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
		log.Fatal("TOKEN is not set — create an app token via the backend and add TOKEN=<value> to .env, then restart: docker compose restart mcp")
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

	transport := os.Getenv("MCP_TRANSPORT")
	if transport == "" {
		stat, err := os.Stdin.Stat()
		if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
			transport = "stdio"
		} else {
			transport = "http"
		}
	}

	switch transport {
	case "stdio":
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	default:
		port := os.Getenv("PORT")
		if port == "" {
			port = "8081"
		}
		handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			return server
		}, nil)
		addr := fmt.Sprintf(":%s", port)
		log.Printf("MCP server listening on %s (streamable HTTP)", addr)
		if err := http.ListenAndServe(addr, handler); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}
}
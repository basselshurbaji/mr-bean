package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/basselshurbaji/mr_bean/backend/config"
	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
)

func main() {
	cfg := config.Load()

	r := handler.NewRouter()

	routes := []handler.Route{
		handler.Adapt(handler.NewHealthHandler()),
	}

	for _, route := range routes {
		handler.Register(r, route)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("server listening on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

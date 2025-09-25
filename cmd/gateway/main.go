package main

import (
	"log"
	"net/http"
	"time"

	"grpc-api-gateway/pkg/loadbalancer"
	"grpc-api-gateway/pkg/middleware"
	"grpc-api-gateway/pkg/server"
)

func main() {
	const serverCount = 5

	// Start the backend test servers
	go server.RunServers(serverCount)

	// Give servers time to start
	time.Sleep(2 * time.Second)

	// Create load balancer and endpoints
	lb, ep := loadbalancer.NewLoadBalancer(serverCount, "http://localhost:808")

	// Setup auth middleware
	auth := middleware.NewAuthMiddleware(middleware.AuthConfig{
		Enabled: true,
		APIKeys: map[string]string{
			"test-key": "test-user",
		},
	})

	// Create router and wrap handler with auth middleware
	router := http.NewServeMux()
	router.Handle("/loadBalancer", auth(loadbalancer.MakeHandler(lb, ep)))

	// Create and start server
	server := &http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	log.Printf("Starting load balancer on :8090")
	log.Fatal(server.ListenAndServe())
}

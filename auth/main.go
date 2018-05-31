package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"

	"github.com/lucperkins/colossus/proto/auth"
)

const (
	PORT = 8888
)

type authHandler struct {
	redisClient *redis.Client
	logger      *fluent.Fluent
}

func (h *authHandler) Authenticate(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {
	var authenticated bool

	password := req.Password

	log.Printf("Request received for the password %s", password)

	value, err := h.redisClient.Get("password").Result()

	if err != nil {
		log.Fatalf("Could not fetch value from Redis: %v", err)
	}

	if password == value {
		authenticated = true
	} else {
		authenticated = false
	}

	return &auth.AuthResponse{Authenticated: authenticated}, nil
}

func main() {
	log, err := fluent.New(fluent.Config{
		FluentPort: 80,
		FluentHost: "elasticsearch",
	})

	if err != nil {
		log.Fatalf("Could not set up fluentd logger: %v", err)
	}

	defer logger.Close()

	log.Printf("Starting up the gRPC auth server on localhost:%d", PORT)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "colossus-redis-cluster:6379",
	})

	_, err := redisClient.Ping().Result()

	if err != nil {
		log.Fatalf("Could not connect to Redis cluster: %v", err)
	}

	defer redisClient.Close()

	log.Print("Successfully connected to Redis")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))

	defer listener.Close()

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	authServer := authHandler{
		redisClient: redisClient,
	}

	auth.RegisterAuthServiceServer(server, &authServer)

	server.Serve(listener)
}

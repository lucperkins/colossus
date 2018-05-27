package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/go-redis/redis"
	"google.golang.org/grpc"

	pb "github.com/lucperkins/colossus/proto/auth"
)

const (
	PORT = 8888
)

type authHandler struct {
	redisClient *redis.Client
}

func (h *authHandler) Authenticate(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
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

	return &pb.AuthResponse{Authenticated: authenticated}, nil
}

func main() {
	log.Printf("Starting up the gRPC auth server on localhost:%d", PORT)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis-master.default.svc.cluster.local:6379",
	})

	_, err := redisClient.Ping().Result()

	if err != nil {
		log.Fatalf("Could not connect to Redis cluster: %v", err)
	}

	log.Print("Successfully connected to Redis")

	redisClient.Set("password", "tonydanza", 0)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	authServer := authHandler{
		redisClient: redisClient,
	}

	pb.RegisterAuthServiceServer(server, &authServer)

	server.Serve(listener)
}

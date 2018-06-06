package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/lucperkins/colossus/proto/auth"
)

const (
	PORT = 8888

	PROMETHEUS_PORT = 9092
)

var (
	metricsRegistry = prometheus.NewRegistry()

	grpcMetrics = grpc_prometheus.NewServerMetrics()

	authCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_svc_success",
		Help: "Auth success counter",
	})

	failCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_svc_fail",
		Help: "Auth fail counter",
	})
)

type authHandler struct {
	redisClient *redis.Client
}

func (h *authHandler) Authenticate(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {
	var authenticated bool

	password := req.Password

	log.Printf("Request received for the password %s", password)

	value, err := h.redisClient.Get("password").Result()

	if err != nil {
		log.Fatalf("Could not fetch value from Redis: %v", err)
	}

	authenticated = password == value

	if authenticated {
		log.Printf("Password %s succeeded", password)
		authCounter.Inc()
	} else {
		log.Printf("Password %s failed", password)
		failCounter.Inc()
	}

	return &auth.AuthResponse{Authenticated: authenticated}, nil
}

func main() {
	log.Printf("Starting up the gRPC auth server on localhost:%d", PORT)

	redisClient := redis.NewClient(&redis.Options{
		//Addr: "colossus-redis-cluster:6379",
		Addr: "localhost:6379",
	})

	_, err := redisClient.Ping().Result()

	if err != nil {
		log.Fatalf("Could not connect to Redis cluster: %v", err)
	}

	log.Print("Successfully connected to Redis")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Print("Successfully created TCP listener")

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)

	authServer := authHandler{
		redisClient: redisClient,
	}

	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", PROMETHEUS_PORT),
	}

	auth.RegisterAuthServiceServer(server, &authServer)

	grpcMetrics.InitializeMetrics(server)

	metricsRegistry.MustRegister(grpcMetrics, authCounter, failCounter)

	log.Print("Successfully registered with Prometheus")

	go func() {
		log.Print("Starting up HTTP server for Prometheus metrics collection")

		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Unable to start HTTP server for Prometheus metrics: %v", err)
		}
	}()

	log.Fatal(server.Serve(listener))
}

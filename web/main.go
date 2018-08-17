package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lucperkins/colossus/proto/auth"
	"github.com/lucperkins/colossus/proto/data"
	"github.com/lucperkins/colossus/proto/userinfo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unrolled/render"
	"google.golang.org/grpc"
)

const (
	port = 3000
)

type (
	httpServer struct {
		authClient     auth.AuthServiceClient
		dataClient     data.DataServiceClient
		renderer       *render.Render
		userInfoClient userinfo.UserInfoClient
		httpReqs       *prometheus.CounterVec
	}

	config struct {
		AuthServiceHost     string `env:"AUTH_SERVICE_HOST"`
		AuthServicePort     string `env:"AUTH_SERVICE_PORT"`
		DataServiceHost     string `env:"DATA_SERVICE_HOST"`
		DataServicePort     string `env:"DATA_SERVICE_PORT"`
		UserinfoServiceHost string `env:"USERINFO_SERVICE_HOST"`
		UserinfoServicePort string `env:"USERINFO_SERVICE_PORT"`
	}
)

func (s *httpServer) PrometheusMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		if r.RequestURI != "/metrics" {
			s.httpReqs.WithLabelValues(http.StatusText(ww.Status()), strings.ToLower(r.Method), r.URL.Path).Inc()
		}
	})
}

func (s *httpServer) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()

		password := r.Header.Get("Password")

		if password == "" {
			http.Error(w, "You cannot access this resource", http.StatusUnauthorized)
			return
		}

		log.Printf("Password attempted: %s", password)

		req := &auth.AuthRequest{
			Password: password,
		}
		res, err := s.authClient.Authenticate(ctx, req)

		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		authenticated := res.Authenticated

		if !authenticated {
			http.Error(w, "You cannot access this resource", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *httpServer) handleString(w http.ResponseWriter, r *http.Request) {
	requestString := r.Header.Get("String")

	if requestString == "" {
		http.Error(w, "You must specify a string using the String header", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	s.dataHandler(ctx, requestString, w)
}

func (s *httpServer) handleStream(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	req := &data.EmptyRequest{}

	stream, err := s.dataClient.StreamingGet(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := []string{}

	for {
		value, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		items = append(items, value.Value)
	}

	s.renderer.JSON(w, http.StatusOK, items)

}

func (s *httpServer) handlePut(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	stream, err := s.dataClient.StreamingPut(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := []string{"foo", "bar", "baz"}

	for _, item := range items {
		req := &data.DataRequest{
			Request: item,
		}

		if err := stream.Send(req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	value := map[string]string{
		"value": res.Value,
	}

	s.renderer.JSON(w, http.StatusAccepted, value)
}

func (s *httpServer) dataHandler(ctx context.Context, requestString string, w http.ResponseWriter) {
	req := &data.DataRequest{
		Request: requestString,
	}

	res, err := s.dataClient.Get(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	value := res.Value

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

func (s *httpServer) handleUserInfo(w http.ResponseWriter, r *http.Request) {
	var username string

	if r.Header.Get("Username") == "" {
		username = "NONE"
	} else {
		username = r.Header.Get("Username")
	}

	ctx := r.Context()

	req := &userinfo.UserInfoRequest{
		Username: username,
	}

	res, err := s.userInfoClient.GetUserInfo(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo := res.UserInfo

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(userInfo))
}

func prometheusWebCounter() *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "web_svc_request_info",
			Help:        "HTTP request counter by response code, request method, and request path",
			ConstLabels: prometheus.Labels{"service": "colossus-web"},
		},
		[]string{"code", "method", "path"},
	)
}

func main() {
	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Could not parse environment variables: %v", err)
	}

	authConn, err := grpc.Dial(
		fmt.Sprintf("colossus-%s-svc:%s", cfg.AuthServiceHost, cfg.AuthServicePort), grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	log.Print("Established connection with auth service")

	dataConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.DataServiceHost, cfg.DataServicePort), grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	log.Print("Established connection with data service")

	userInfoConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.UserinfoServiceHost, cfg.UserinfoServicePort), grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	log.Print("Established connection with userinfo service")

	authClient := auth.NewAuthServiceClient(authConn)
	dataClient := data.NewDataServiceClient(dataConn)
	userInfoClient := userinfo.NewUserInfoClient(userInfoConn)

	r := chi.NewRouter()

	renderer := render.New(render.Options{})

	httpReqs := prometheusWebCounter()

	if err := prometheus.Register(httpReqs); err != nil {
		log.Fatalf("Could not register Prometheus httpReqs counter vec: %v", err)
	}

	server := httpServer{
		authClient:     authClient,
		dataClient:     dataClient,
		renderer:       renderer,
		userInfoClient: userInfoClient,
		httpReqs:       httpReqs,
	}

	log.Print("Using the following middleware: Prometheus metrics, authentication")

	// The Prometheus metrics middleware
	r.Use(server.PrometheusMetrics)

	// The authentication layer
	r.Use(server.authenticate)

	r.Post("/string", server.handleString)

	r.Get("/stream", server.handleStream)

	r.Put("/stream", server.handlePut)

	r.Get("/user", server.handleUserInfo)

	// The Prometheus metrics handler
	r.Handle("/metrics", prometheus.Handler())

	log.Printf("Now starting the server on port %d...", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

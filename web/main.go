package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/lucperkins/colossus/proto/auth"
	"github.com/lucperkins/colossus/proto/data"
	"google.golang.org/grpc"
)

const (
	PORT              = 3000
	AUTH_SERVICE_PORT = 8888
	DATA_SERVICE_PORT = 1111
)

type HttpServer struct {
	authClient auth.AuthServiceClient
	dataClient data.DataServiceClient
}

func (s *HttpServer) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func (s *HttpServer) handleString(w http.ResponseWriter, r *http.Request) {
	requestString := r.Header.Get("String")

	if requestString == "" {
		http.Error(w, "You must specify a string using the String header", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	s.dataHandler(ctx, requestString, w)
}

func (s *HttpServer) dataHandler(ctx context.Context, requestString string, w http.ResponseWriter) {
	req := &data.DataRequest{
		Request: requestString,
	}

	res, err := s.dataClient.Get(ctx, req)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	value := res.Value

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

func main() {
	authConn, err := grpc.Dial(
		fmt.Sprintf("colossus-auth-svc:%d", AUTH_SERVICE_PORT), grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	log.Print("Established connection with auth service")

	dataConn, err := grpc.Dial(
		fmt.Sprintf("colossus-data-svc:%d", DATA_SERVICE_PORT), grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	log.Print("Established connection with data service")

	authClient := auth.NewAuthServiceClient(authConn)
	dataClient := data.NewDataServiceClient(dataConn)

	r := chi.NewRouter()

	server := HttpServer{
		authClient: authClient,
		dataClient: dataClient,
	}

	log.Print("Using the following middleware: authentication")

	r.Use(server.authenticate)

	r.Post("/string", server.handleString)

	log.Printf("Now starting the server on port %d...", PORT)

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
}

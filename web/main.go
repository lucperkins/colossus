package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/lucperkins/colossus/proto/auth"
	"github.com/lucperkins/colossus/proto/data"
	"google.golang.org/grpc"
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

func (s *HttpServer) Get(w http.ResponseWriter, r *http.Request) {
	word := r.Header.Get("Word")

	if word == "" {
		http.Error(w, "You must specify a word using the Word header", http.StatusBadRequest)
		return
	}

	req := &data.DataRequest{
		Key: word,
	}

	ctx := r.Context()

	res, err := s.dataClient.Get(ctx, req)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.Value))
}

func main() {
	port := 3000

	authConn, err := grpc.Dial("colossus-auth-svc.default.svc.cluster.local:8888", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	log.Print("Established connection with auth service")

	dataConn, err := grpc.Dial("colossus-data-svc.default.svc.cluster.local:1111", grpc.WithInsecure())

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

	r.Get("/", server.Get)

	log.Printf("Now starting the server on port %d...", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

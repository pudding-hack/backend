package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/pudding-hack/backend/lib"
	"github.com/rs/cors"
)

type handler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	GetItemHistoryPaginate(w http.ResponseWriter, r *http.Request)
	InboundItem(w http.ResponseWriter, r *http.Request)
	OutboundItem(w http.ResponseWriter, r *http.Request)
	DetectLabels(w http.ResponseWriter, r *http.Request)
}

func Run(ctx context.Context, cfg lib.Config, requestHandler handler) error {
	router := mux.NewRouter()
	inventoryMiddleware := lib.NewAuthMiddleware(&cfg)
	router.Use(inventoryMiddleware.ValidateCurrentUser)

	// grouping handler
	api := router.PathPrefix("/api").Subrouter()
	inventory := api.PathPrefix("/inventory").Subrouter()

	// inventory handler
	inventory.HandleFunc("/", requestHandler.GetAll).Methods(http.MethodGet)
	inventory.HandleFunc("/detail", requestHandler.GetById).Methods(http.MethodGet)
	inventory.HandleFunc("/create", requestHandler.Create).Methods(http.MethodPost)

	// history handler
	inventory.HandleFunc("/history", requestHandler.GetItemHistoryPaginate).Methods(http.MethodGet)

	// inbound outbound handler
	inventory.HandleFunc("/inbound", requestHandler.InboundItem).Methods(http.MethodPost)
	inventory.HandleFunc("/outbound", requestHandler.OutboundItem).Methods(http.MethodPost)

	// rekognition handler
	inventory.HandleFunc("/rekognition", requestHandler.DetectLabels).Methods(http.MethodPost)

	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"POST", "GET", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Mode"},
		MaxAge:             60, // 1 minutes
		AllowCredentials:   true,
		OptionsPassthrough: false,
		Debug:              false,
	})

	httpHandler := c.Handler(router)

	err := startServer(ctx, httpHandler, cfg)
	if err != nil {
		return err
	}

	return nil
}

func startServer(ctx context.Context, httpHandler http.Handler, cfg lib.Config) error {
	errChan := make(chan error)

	go func() {
		errChan <- startHTTP(ctx, httpHandler, cfg)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func startHTTP(ctx context.Context, httpHandler http.Handler, cfg lib.Config) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.HTTPPort),
		Handler: httpHandler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server: ", err)
		}
	}()

	log.Printf("%s is starting at port: %d", cfg.App.Name, cfg.App.HTTPPort)
	interruption := make(chan os.Signal, 1)
	signal.Notify(interruption, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-interruption

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown: %s", err)
		return err
	}

	log.Println("server is shutting down")
	return nil
}

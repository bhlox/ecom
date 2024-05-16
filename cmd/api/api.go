package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/bhlox/ecom/internal/db"
	"github.com/bhlox/ecom/internal/health"
	"github.com/bhlox/ecom/internal/services/checkout"
	"github.com/bhlox/ecom/internal/services/order"
	"github.com/bhlox/ecom/internal/services/product"
	"github.com/bhlox/ecom/internal/services/user"
)

type APIServer struct {
	address string
	db      *sql.DB
}

func NewAPIServer(address string, db *sql.DB) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}
}

// this func handles routing for all routes. for speciific routes check the registerRoutes method for each handler
func (s *APIServer) Run() error {
	router := http.NewServeMux()

	healthRouter := http.NewServeMux()
	healthHandler := health.NewHandler()
	healthHandler.RegisterRoutes(healthRouter)

	// Add the /health prefix to the healthRouter
	router.Handle("/health/", http.StripPrefix("/health", healthRouter))

	v1Router := http.NewServeMux()
	v1Router.Handle("/v1/", http.StripPrefix("/v1", router))

	userHandler := user.NewHandler(db.New(s.db))
	userHandler.RegisterRoutes(v1Router)

	productHandler := product.NewHandler(db.New(s.db))
	productHandler.RegisterRoutes(v1Router)

	checkoutHandler := checkout.NewHandler(db.New(s.db))
	checkoutHandler.RegisterRoutes(v1Router)

	ordersHandler := order.NewHandler(db.New(s.db))
	ordersHandler.RegisterRoutes(v1Router)

	// Register the v1 router to the main router
	router.Handle("/v1/", http.StripPrefix("/v1", v1Router))

	//  // Set security headers
	//  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//     w.Header().Set("Content-Security-Policy", "default-src 'self'")
	//     w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	//     w.Header().Set("X-Content-Type-Options", "nosniff")
	//     w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	//     w.WriteHeader(http.StatusOK)
	// })

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: time.Second * 30,
	}

	log.Printf("🚀 server listening at localhost:%s", "8080")
	return server.ListenAndServe()
}

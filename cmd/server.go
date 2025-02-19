package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"vibhordubey333/loan-service/config"
	"vibhordubey333/loan-service/internal/repository"
	"vibhordubey333/loan-service/internal/handler"
	"vibhordubey333/loan-service/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Could not establish a connection to the database: %v", err)
	} else {
		log.Println("Successfully connected to the database")
	}

	loanRepo := repository.NewLoanRepository(db)

	loanService := service.NewLoanService(loanRepo)

	// Handler layer
	loanHandler := handler.NewLoanHandler(loanService)
	//Router setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/loans", func(r chi.Router) {
			r.Post("/", loanHandler.CreateLoan)
			r.Get("/{id}", loanHandler.GetLoan)
		})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Graceful shutdown
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	log.Printf("Server is running on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", cfg.Port, err)
	}

	<-done
	log.Println("Server stopped")
}

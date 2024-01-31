package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"wallet/api"
	"wallet/repository"
	"wallet/service"
)

const (
	DEFAULT_RDS_HOST     = "localhost"
	DEFAULT_RDS_USER     = "postgres"
	DEFAULT_RDS_PASSWORD = "postgres"
)

func main() {
	r := mux.NewRouter()

	db, err := connectToDB()
	if err != nil {
		slog.Error("Could not connect to db", err)
		return
	}

	repo, err := repository.NewRepository(db)
	if err != nil {
		slog.Error("Could not init repo", err)
		return
	}

	api := api.NewAPI(service.NewService(repo))

	r.HandleFunc("/wallet/{wallet_id}", panicHandler((api.GetWallet))).Methods("GET")
	r.HandleFunc("/wallet/{wallet_id}/add", panicHandler((api.AddToWallet))).Methods("POST")
	r.HandleFunc("/wallet/{wallet_id}/subtract", panicHandler((api.SubtractFromWallet))).Methods("POST")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	slog.Info("listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func panicHandler(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			panicValue := recover()

			if panicValue == nil {
				return
			}

			err, ok := panicValue.(error)
			if !ok {
				err = errors.New("got a non-error panic")
			}

			slog.Error("recovered from panic", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}()

		next(w, r)
	})
}

func connectToDB() (*sql.DB, error) {
	var RDS_HOST string = os.Getenv("RDS_HOST")
	if RDS_HOST == "" {
		RDS_HOST = DEFAULT_RDS_HOST
	}
	var RDS_USER string = os.Getenv("RDS_USER")
	if RDS_USER == "" {
		RDS_USER = DEFAULT_RDS_USER
	}
	var RDS_PASSWORD string = os.Getenv("RDS_PASSWORD")
	if RDS_PASSWORD == "" {
		RDS_PASSWORD = DEFAULT_RDS_PASSWORD
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s?sslmode=disable", RDS_USER, RDS_PASSWORD, RDS_HOST)
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

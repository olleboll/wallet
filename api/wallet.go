package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	models "wallet/models"

	"github.com/gorilla/mux"
)

type updateBalanceRequest struct {
	Amount int32 `json:"amount"`
}

type walletService interface {
	GetWallet(ctx context.Context, WalletID string) (*models.Wallet, error)
	AddToWallet(ctx context.Context, WalletID string, amount int32) error
	SubtractFromWallet(ctx context.Context, WalletID string, amount int32) error
}

type walletAPI struct {
	walletService walletService
}

func NewAPI(service walletService) *walletAPI {
	return &walletAPI{walletService: service}
}

func (api *walletAPI) GetWallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID := vars["wallet_id"]

	wallet, err := api.walletService.GetWallet(r.Context(), walletID)
	if err != nil {
		slog.Error(fmt.Sprintf("Received unexpected error:\n%+v\n", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if wallet.WalletID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	responseData, err := json.Marshal(wallet)
	if err != nil {
		slog.Error(fmt.Sprintf("Received unexpected error:\n%+v\n", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(responseData)
	if err != nil {
		slog.Error(fmt.Sprintf("Received unexpected error:\n%+v\n", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (api *walletAPI) AddToWallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID := vars["wallet_id"]

	var request updateBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.Error(fmt.Sprintf("Received unexpected error:\n%+v\n", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := api.walletService.AddToWallet(r.Context(), walletID, request.Amount)
	if err != nil {
		slog.Error(fmt.Sprintf("Received unexpected error:\n%+v\n", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *walletAPI) SubtractFromWallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID := vars["wallet_id"]

	var request updateBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return
	}

	err := api.walletService.SubtractFromWallet(r.Context(), walletID, request.Amount)
	if err != nil {
		slog.Error(fmt.Sprintf("Received unexpected error:\n%+v\n", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

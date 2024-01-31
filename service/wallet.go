package service

import (
	"context"

	models "wallet/models"
)

type walletRepository interface {
	GetWallet(ctx context.Context, walletID string) (*models.Wallet, error)
	UpsertWalletBalance(ctx context.Context, walletID string, amount int32) error
}

type service struct {
	repo walletRepository
}

func (svc *service) GetWallet(ctx context.Context, walletID string) (*models.Wallet, error) {
	return svc.repo.GetWallet(ctx, walletID)
}

func (svc *service) AddToWallet(ctx context.Context, walletID string, amount int32) error {
	return svc.repo.UpsertWalletBalance(ctx, walletID, amount)
}

func (svc *service) SubtractFromWallet(ctx context.Context, walletID string, amount int32) error {
	return svc.repo.UpsertWalletBalance(ctx, walletID, -amount)
}

func NewService(repo walletRepository) *service {
	return &service{repo: repo}
}

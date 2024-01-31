package repository

import (
	"context"
	"database/sql"
	"errors"

	"wallet/models"

	_ "github.com/lib/pq"
)

type walletRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) (*walletRepository, error) {
	// For simplicity!
	if err := createTable(db); err != nil {
		return nil, err
	}

	return &walletRepository{db: db}, nil
}

func (repo *walletRepository) GetWallet(ctx context.Context, walletID string) (*models.Wallet, error) {
	row := repo.db.QueryRow(`SELECT id, balance FROM "wallets" WHERE "id" = $1`, walletID)

	var wallet models.Wallet
	err := row.Scan(&wallet.WalletID, &wallet.Balance)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &wallet, nil
}

func (repo *walletRepository) UpsertWalletBalance(ctx context.Context, walletID string, amount int32) error {
	wallet, err := repo.GetWallet(ctx, walletID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// If the wallet does not exists, we create it
	sql := `UPDATE wallets SET balance = balance + $1 WHERE id = $2;`
	if wallet.WalletID == "" {
		sql = `INSERT INTO wallets (balance, id) VALUES ($1, $2)`
	}

	_, err = repo.db.Exec(sql, amount, walletID)
	return err
}

func createTable(db *sql.DB) error {
	sql := `CREATE TABLE IF NOT EXISTS wallets (
		id text PRIMARY KEY,
		balance int
	)`
	_, err := db.Exec(sql)
	return err
}

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"wallet/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type updateWalletRequest struct {
	Amount int32 `json:"amount"`
}

func TestWallet(t *testing.T) {
	client := http.Client{}

	walletID := uuid.New()

	// Create and add balance
	initialWalletRequest := updateWalletRequest{
		Amount: 100,
	}

	requestBytes, err := json.Marshal(initialWalletRequest)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:8080/wallet/%s/add", walletID), bytes.NewReader(requestBytes))
	require.NoError(t, err)

	response, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Read balance
	var walletResponse models.Wallet
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/wallet/%s", walletID), nil)
	require.NoError(t, err)

	response, err = client.Do(req)
	require.NoError(t, err)

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &walletResponse)
	require.NoError(t, err)

	assert.Equal(t, walletID.String(), walletResponse.WalletID)
	assert.Equal(t, initialWalletRequest.Amount, walletResponse.Balance)

	// Subtract balance
	buyUsingWalletRequest := updateWalletRequest{
		Amount: 42,
	}

	requestBytes, err = json.Marshal(buyUsingWalletRequest)
	require.NoError(t, err)

	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:8080/wallet/%s/subtract", walletID), bytes.NewReader(requestBytes))
	require.NoError(t, err)

	response, err = client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Read again
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/wallet/%s", walletID), nil)
	require.NoError(t, err)

	response, err = client.Do(req)
	require.NoError(t, err)

	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &walletResponse)
	require.NoError(t, err)

	assert.Equal(t, initialWalletRequest.Amount-buyUsingWalletRequest.Amount, walletResponse.Balance)
}

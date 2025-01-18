package tx_fees

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Tx struct {
	Hash      string `json:"hash"`
	BlockHash string `json:"blockHash"`
	BlockNum  string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	GasUsed   string `json:"gasUsed"`
	GasPrice  string `json:"gasPrice"`
	GasCost   string `json:"gasCost"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	TimeStamp string `json:"timeStamp"`
}

// TxListResponse represents the API response for transactions.
type TxListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []Tx   `json:"result"`
}

const (
	apiURL = "https://api.arbiscan.io/api"
)

// getTransactions fetches all transactions for the given address.
func getAllArbitrumTransactions(db *sql.DB, address string) (error) {

	apiKey := os.Getenv("ARB_API_KEY")
	if apiKey == "" {
		log.Logger.Error().Msg("ARB_API_KEY environment variable is not set")
	}

	url := fmt.Sprintf("%s?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s", apiURL, address, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to read http response")
		return fmt.Errorf("failed to read response: %v", err)
	}

	var txResponse TxListResponse
	if err := json.Unmarshal(body, &txResponse); err != nil {
		log.Logger.Error().Err(err).Msg("Failed to parse Arbitrum txs JSON from http response")
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	//parse txResponse.Status to int
	status, err := strconv.Atoi(txResponse.Status)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("Error parsing Status for address %s", address)
		return fmt.Errorf("Error parsing Status for address %s", address)
	}

	if status == 0 {
		log.Logger.Error().Msgf("API error: %s", txResponse.Message)
		return fmt.Errorf("API error: %s", txResponse.Message)
	}

	
	for _, tx := range txResponse.Result {
		// Calculate the gas cost
		gasUsed, err := strconv.ParseFloat(tx.GasUsed, 64)
		if err != nil {
			log.Logger.Error().Err(err).Msgf("Error parsing GasUsed for transaction %s", tx.Hash)
			continue
		}
		gasPrice, err := strconv.ParseFloat(tx.GasPrice, 64)
		if err != nil {
			log.Logger.Error().Err(err).Msgf("Error parsing GasPrice for transaction %s", tx.Hash)
			continue
		}
		tx.GasCost = fmt.Sprintf("%f", gasUsed * gasPrice)

		// Insert the transaction into the database
		err = insertTransaction(db, tx)
		if err != nil {
			log.Logger.Error().Err(err).Msgf("Error inserting transaction %s", tx.Hash)
		}
	}
	return nil
}

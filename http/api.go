package http

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rs/zerolog/log"
)

func QueryHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, solverAddress string) {

	if db == nil {
        http.Error(w, "Database connection is not initialized", http.StatusInternalServerError)
        return
    }

	// Parse query parameters from the URL
	chainID := r.URL.Query().Get("chain_id")
	if chainID == "" {
		http.Error(w, "Missing chain_id query parameter", http.StatusBadRequest)
		return
	}

	// Query the database using the chain_id
	var totalOrderCount int
	row := db.QueryRow("SELECT COUNT(*) FROM tx_data WHERE source_domain = ?", chainID)
	err := row.Scan(&totalOrderCount)
	if err != nil {
		log.Logger.Error().Msg("Failed to execute totalOrderCount query")
	}

	// Query the database using the chain_id
	var totalOrderFilled int
	row = db.QueryRow("SELECT COUNT(*) FROM tx_data WHERE source_domain = ? AND filler = ?", chainID, solverAddress)
	err = row.Scan(&totalOrderFilled)
	if err != nil {
		log.Logger.Error().Msg("Failed to execute totalOrderFilled query")
	}

	// Calculate the success rate
	rate := float32(totalOrderFilled) / float32(totalOrderCount)

	// Query the database to get the total solver revenue
	var totalRevenue int
	err = db.QueryRow("SELECT SUM(solver_revenue) FROM tx_data WHERE source_domain = ? AND filler = ?", chainID, solverAddress).Scan(&totalRevenue)
    if err != nil {
        if err == sql.ErrNoRows {
			log.Logger.Error().Msg("No rows found")
            return
        }
        return
    }
	revenue := float32(totalRevenue) / 1000000

	var totalOsmoFees float32
	err = db.QueryRow("SELECT SUM(fee_amount) FROM tx_data WHERE source_domain = ? AND filler = ?", chainID, solverAddress).Scan(&totalOsmoFees)
    if err != nil {
        if err == sql.ErrNoRows {
			log.Logger.Error().Msg("No rows found")
            return
        }
        return
    }
	osmoFees := float32(totalOsmoFees) / 1000000

	// Osmo fees in USDC
	var osmoPrice float32
	err = db.QueryRow("SELECT price_usd FROM eth_prices WHERE token_denom = 'osmosis'").Scan(&osmoPrice)
	if err != nil {
		log.Logger.Error().Msg("Failed to execute osmoPrice query")
	}
	totalOsmoFeesUSD := osmoFees * osmoPrice

	// Calculate all tx fees from Arbitrum
	var totalArbitrumTxFees float32
	err = db.QueryRow("SELECT SUM(gas_cost) FROM arbitrum_txs").Scan(&totalArbitrumTxFees)
	if err != nil {
		log.Logger.Error().Msg("Failed to execute totalArbitrumTxFees query")
	}
	totalArbitrumTxFees = totalArbitrumTxFees / 1000000000000000000

	// Arbitrum Eth fees in USDC
	var ethPrice float32
	err = db.QueryRow("SELECT price_usd FROM eth_prices WHERE token_denom = 'arbitrum'").Scan(&ethPrice)
	if err != nil {
		log.Logger.Error().Msg("Failed to execute osmoPrice query")
	}
	totalArbitrumTxFeesUSD := totalArbitrumTxFees * ethPrice


	resp := ApiResponse{
		SrcChain:     chainID,
		DstChain:     "osmosis-1",
		TotalOrders:  totalOrderCount,
		TotalFilled:  totalOrderFilled,
		TotalRevenueUSDC: revenue,
		TotalOsmosisTxFeesUSD: totalOsmoFeesUSD,
		TotalArbitrumTxFeesUSD: totalArbitrumTxFeesUSD,
		TotalArbitrumTxFeesETH: totalArbitrumTxFees,
		TotalOsmosisTxFeesOSMO: osmoFees,
		SuccessRate:  rate,
	}

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

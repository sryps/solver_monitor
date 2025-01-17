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


	rate := float32(totalOrderFilled) / float32(totalOrderCount)
	revenue := float32(totalRevenue) / 1000000

	resp := ApiResponse{
		SrcChain:     chainID,
		DstChain:     "osmosis-1",
		TotalOrders:  totalOrderCount,
		TotalFilled:  totalOrderFilled,
		TotalRevenueUSDC: revenue,
		SuccessRate:  rate,
	}

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

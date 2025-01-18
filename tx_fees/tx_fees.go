package tx_fees

import (
	"database/sql"

	"github.com/rs/zerolog/log"

	_ "github.com/mattn/go-sqlite3"
)


func CollectAllTxFees (db *sql.DB, arbitrumSolverAddress *string) error {

	log.Logger.Info().Msg("Collecting all transaction fees")
	// Get all transactions for the arbitrum address
	err := getAllArbitrumTransactions(db, *arbitrumSolverAddress)
	if err != nil {
		return err
	}
	log.Logger.Info().Msg("Completed colecting all transaction fees")
	return nil
}
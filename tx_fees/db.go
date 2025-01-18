package tx_fees

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"

	_ "github.com/mattn/go-sqlite3"
)

func InitArbitrumTxFeesDB(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS arbitrum_txs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			hash TEXT UNIQUE NOT NULL,
			from_address TEXT NOT NULL,
			to_address TEXT NOT NULL,
			gas_used TEXT NOT NULL,
			gas_price TEXT NOT NULL,
			gas_cost TEXT NOT NULL,
			value TEXT NOT NULL,
			timestamp TEXT NOT NULL
		);`
	_, err := db.Exec(query)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create arbitrum_txs table")
	}
}

func insertTransaction(db *sql.DB, tx Tx) error {
	// Check if the transaction hash already exists
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM arbitrum_txs WHERE hash = ?", tx.Hash).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check transaction existence: %v", err)
	}

	if exists > 0 {
		// Transaction already recorded, skip insertion
		return nil
	}

	log.Logger.Info().Msgf("Inserting Arbitrum TX %s fee data into DB", tx.Hash)

	// Insert the transaction into the database
	_, err = db.Exec(`
		INSERT INTO arbitrum_txs (hash, from_address, to_address, gas_used, gas_price, gas_cost, value, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		tx.Hash, tx.From, tx.To, tx.GasUsed, tx.GasPrice, tx.GasCost, tx.Value, tx.TimeStamp,
	)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}

	return nil
}
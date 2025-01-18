package tx_fees

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type CoinGeckoAPIResponseETHarb struct {
	Arbitrum struct {
		USD float64 `json:"usd"`
	} `json:"arbitrum"`
}

type CoinGeckoAPIResponseOSMO struct {
	Osmosis struct {
		USD float64 `json:"usd"`
	} `json:"osmosis"`
}

func GetCoingeckoPrice(db *sql.DB) (error) {
	denoms := []string{"arbitrum", "osmosis"}
	for _, denom := range denoms {
		EthPrice, err := fetchPrice(denom)
		if err != nil {
			log.Logger.Error().Err(err).Msgf("Failed to fetch ETH price for %s", denom)
			return err
		}

		err = storePrice(db, EthPrice, denom)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to store ETH price in database")
			return err
		}
		log.Logger.Info().Msgf("Fetched and stored ETH price for %s", denom)
	}
	return nil
}

func fetchPrice(denom string) (float64, error) {

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", denom)

	// Make HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	switch denom {
	case "arbitrum":
		var apiResponse CoinGeckoAPIResponseETHarb
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			return 0, err
		}
		return apiResponse.Arbitrum.USD, nil
	case "osmosis":
		var apiResponse CoinGeckoAPIResponseOSMO
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			return 0, err
		}
		return apiResponse.Osmosis.USD, nil
	default:
		return 0, fmt.Errorf("unsupported denom: %s", denom)
	}
}
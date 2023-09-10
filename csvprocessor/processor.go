package csvprocessor

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	baseURL = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
)

func ProcessCSVStockFile(sKey string) string {
	hydratedURL := fmt.Sprintf(baseURL, url.QueryEscape(sKey))
	log.Info().Msgf("hydratedURL: %s", hydratedURL)
	// call the url , get the csv and process it
	response, err := http.Get(hydratedURL)
	if err != nil {
		log.Err(err).Msg("error calling stock service")
		return "Service is not available"
	}
	switch response.StatusCode {
	case http.StatusOK:
		return readCSVStockFile(response.Body)
	case http.StatusNotFound:
		return "Stock not found"
	default:
		return "Service is not available"
	}

}

// readCSVStockFile csv contents from response.Body
func readCSVStockFile(contents io.ReadCloser) string {
	content, err := csv.NewReader(contents).ReadAll()
	log.Info().Msgf("content: %v", content)
	if err != nil {
		log.Err(err).Msg("error reading csv file")
		return "Stock service CSV error"
	}
	//contents:
	// Symbol,Date,Time,Open,High,Low,Close,Volume
	// AAPL.US,2023-09-08,22:00:10,178.35,180.239,177.79,178.18,65602066
	symbol := content[1][0]
	close := content[1][6]
	log.Info().Msgf("symbol: %s, close: %s", symbol, close)
	_, err = strconv.ParseFloat(close, 64)
	if err != nil {
		return fmt.Sprintf("%s quote is not available", strings.ToUpper(symbol))
	}
	return fmt.Sprintf("%s quote is $%s per share", strings.ToUpper(symbol), close)

}

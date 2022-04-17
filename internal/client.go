package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type OptionType int

const (
	Call OptionType = iota
	Put
)

type Option struct {
	UnderlyingTicker string
	UnderlyingPrice  float32
	ExpirationDate   string
	StrikePrice      float32
	BidPremium       float32
	AskPremium       float32
	OptionType       OptionType
}

type ResultPrices struct {
	BidPrice float32 `json:"b"`
	AskPrice float32 `json:"a"`
}

type Result struct {
	ResultOptions map[string]map[string]map[string]ResultPrices `json:"options"`
}

var HttpClient http.Client = http.Client{
	Timeout: time.Second * 2, // Timeout after 2 seconds
}

func GetPriceForTicker(ticker string) float32 {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/download/%s", ticker)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	res, getErr := HttpClient.Do(req)

	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	tokens := strings.Split(string(body), ",")

	// tokens[10] is Close on the CSV
	result, parseErr := strconv.ParseFloat(tokens[10], 32)
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	return float32(result)
}

func FormatOption(contract Option) string {
	formattedType := "C"
	if contract.OptionType == Put {
		formattedType = "P"
	}
	return fmt.Sprintf("%s (%.2f), Exp. %s, %.2f%s, %.2f - %.2f", contract.UnderlyingTicker, contract.UnderlyingPrice, contract.ExpirationDate, contract.StrikePrice, formattedType, contract.BidPremium, contract.AskPremium)
}

func GetOptionsForTicker(ticker string) []Option {
	url := fmt.Sprintf("https://www.optionsprofitcalculator.com/ajax/getOptions?stock=%s&reqId=1", ticker)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	res, getErr := HttpClient.Do(req)

	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var result Result
	json.Unmarshal(body, &result)

	allContracts := make([]Option, 0)

	underlyingPrice := GetPriceForTicker(ticker)

	for expirationDate, chain := range result.ResultOptions {
		for side, contracts := range chain {
			for strike, prices := range contracts {
				contract := Option{}
				contract.UnderlyingTicker = ticker
				contract.UnderlyingPrice = underlyingPrice
				contract.ExpirationDate = expirationDate
				parsedStrike, parseErr := strconv.ParseFloat(strike, 32)
				if parseErr != nil {
					log.Fatal(parseErr)
				}
				contract.StrikePrice = float32(parsedStrike)
				contract.BidPremium = prices.BidPrice
				contract.AskPremium = prices.AskPrice
				if side == "c" {
					contract.OptionType = Call
				}
				if side == "p" {
					contract.OptionType = Put
				}

				allContracts = append(allContracts, contract)
			}
		}
	}

	return allContracts
}

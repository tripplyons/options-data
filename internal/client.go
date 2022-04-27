package internal

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
	LastPremium      float32
	OpenInterest     int
	Volume           int
	OptionType       OptionType
	TimeSeen         int
}

type ResultPrices struct {
	BidPremium   float32 `json:"b"`
	AskPremium   float32 `json:"a"`
	LastPremium  float32 `json:"l"`
	OpenInterest int     `json:"oi"`
	Volume       int     `json:"v"`
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
	formattedType := "c"
	if contract.OptionType == Put {
		formattedType = "p"
	}
	items := []string{
		fmt.Sprintf("%d", contract.TimeSeen),
		contract.UnderlyingTicker,
		fmt.Sprintf("%.2f", contract.UnderlyingPrice),
		contract.ExpirationDate,
		fmt.Sprintf("%.2f", contract.StrikePrice),
		formattedType,
		fmt.Sprintf("%.2f", contract.BidPremium),
		fmt.Sprintf("%.2f", contract.AskPremium),
		fmt.Sprintf("%.2f", contract.LastPremium),
		fmt.Sprintf("%d", contract.OpenInterest),
		fmt.Sprintf("%d", contract.Volume),
	}
	return strings.Join(items, ",")
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

	now := time.Now()

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
				contract.BidPremium = prices.BidPremium
				contract.AskPremium = prices.AskPremium
				contract.LastPremium = prices.LastPremium
				contract.OpenInterest = prices.OpenInterest
				contract.Volume = prices.Volume
				if side == "c" {
					contract.OptionType = Call
				}
				if side == "p" {
					contract.OptionType = Put
				}
				contract.TimeSeen = int(now.Unix())

				allContracts = append(allContracts, contract)
			}
		}
	}

	return allContracts
}

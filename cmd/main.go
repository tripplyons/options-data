package main

import (
	"fmt"
	"os"
	"strings"

	inter "github.com/tripplyons/options-data/internal"
)

func main() {
	tickers := strings.Split(os.Args[1], ",")

	for _, ticker := range tickers {
		options := inter.GetOptionsForTicker(ticker)

		for _, option := range options {
			fmt.Println(inter.FormatOption(option))
		}
	}
}

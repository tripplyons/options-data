package main

import (
	"fmt"

	inter "github.com/tripplyons/options-data/internal"
)

func main() {
	fmt.Printf("%f\n", inter.GetPriceForTicker("SPY"))
	options := inter.GetOptionsForTicker("SPY")

	for _, option := range options[:10] {
		fmt.Println(inter.FormatOption(option))
	}
}

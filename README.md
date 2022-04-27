# Options Data
Downloads option prices from an API from [optionsprofitcalculator.com](https://optionsprofitcalculator.com/) and outputs them in CSV format.

Usage example:
```sh
echo "timestamp,ticker,underlyingPrice,expiration,strike,type,bid,ask,last,openInterest,volume" > output.csv
go run cmd/main.go SPY,QQQ >> output.csv
```

# ![stock_data](./assets/application.png)

[![GitHub](https://img.shields.io/github/license/bzdanowicz/stock_data?style=for-the-badge)](https://github.com/bzdanowicz/stock_data/blob/master/LICENSE)
[![GitHub last commit](https://img.shields.io/github/last-commit/bzdanowicz/stock_data?style=for-the-badge)](https://github.com/bzdanowicz/stock_data/commits/master)
![Lines of code](https://img.shields.io/tokei/lines/github/bzdanowicz/stock_data?style=for-the-badge)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/bzdanowicz/stock_data?style=for-the-badge)

## Table of contents
* [General info](#general-info)
* [Technologies](#technologies)
* [Setup](#setup)
* [Configuration](#configuration)

## General info
Terminal application designed to provide an easy way for users to get up-to-date, basic price information of all the major stocks trading in the US. Data is pulled from the Finnhub API.
	
## Technologies
Project is created with:
* Golang
* [termui](https://github.com/gizak/termui)

## Setup

1. Clone the repository
2. Add API key and stock symbols to `config.json` file
3. Run application with following command:
```
$ go run cmd/main.go
```

## Configuration

The configuration file is located at `~/config.json`
User should provide own API Key (obtained from [Finnhub](https://finnhub.io/)), list of requested stock ticker symbols and base currency to calculate rates.

```
{
    "apiKey": "apiKey",
    "quotes": [
        "AAPL", "MSFT", "AMZN", "GOOG", "APVO"
    ],
    "base": "PLN"
}
```


package internal

import (
	"github.com/bzdanowicz/stock_data/finnhub"
)

type QuoteData map[string]finnhub.Quote
type RatesData finnhub.Rates
type CandleData finnhub.Candle

type Data struct {
	Quotes QuoteData
	Rates  RatesData
	Candle finnhub.Candle
}

func (data *Data) Initialize(quotes []string, base string) {
	data.Quotes = make(QuoteData)
	for _, q := range quotes {
		data.Quotes[q] = finnhub.Quote{}
	}
	data.Rates.Base = base
}

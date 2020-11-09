package internal

import (
	"github.com/bzdanowicz/stock_data/finnhub"
)

type QuoteData map[string]finnhub.Quote
type RatesData finnhub.Rates

type Data struct {
	Quotes QuoteData
	Rates  RatesData
}

func (data *Data) Initialize(quotes []string, base string) {
	data.Quotes = make(QuoteData)
	data.Rates.Base = base
	for _, q := range quotes {
		data.Quotes[q] = finnhub.Quote{}
	}
}

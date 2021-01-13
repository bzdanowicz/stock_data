package finnhub

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

type Candle struct {
	CurrentPrice []float64 `json:"c"`
	LowPrice     []float32 `json:"l"`
	HighPrice    []float32 `json:"h"`
	OpenPrice    []float32 `json:"o"`
	Volume       []float32 `json:"v"`
	Timestamps   []int64   `json:"t"`
	Status       string    `json:"s"`
}

func (c *Client) GetCandle(symbol string, from time.Time, to time.Time, resolution string) (*Candle, error) {
	query, _ := url.Parse(c.baseURL.String())
	query.Path += "/stock/candle"
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("resolution", resolution)
	params.Add("from", strconv.FormatInt(from.Unix(), 10))
	params.Add("to", strconv.FormatInt(to.Unix(), 10))
	query.RawQuery = params.Encode()

	res, err := c.Get(query.String())
	if err != nil {
		return &Candle{}, err
	}

	candle := &Candle{}
	err = json.NewDecoder(res.Body).Decode(candle)
	return candle, err
}

type CandleTask struct {
	Symbol     string
	From       time.Time
	To         time.Time
	Resolution string
	Executor   func(symbol string, from time.Time, to time.Time, resolution string) (*Candle, error)
}

func (c *CandleTask) Execute() (interface{}, error) {
	return c.Executor(c.Symbol, c.From, c.To, c.Resolution)
}

func (q *CandleTask) GetParameter() string {
	return q.Symbol
}

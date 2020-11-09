package finnhub

import (
	"encoding/json"
	"net/url"
)

type RatesValues struct {
	CHF float32 `json:"CHF"`
	EUR float32 `json:"EUR"`
	GBP float32 `json:"GBP"`
	JPY float32 `json:"JPY"`
	PLN float32 `json:"PLN"`
	USD float32 `json:"USD"`
}

type Rates struct {
	Base   string      `json:"base"`
	Values RatesValues `json:"quote"`
}

func (c *Client) GetRates(base string) (*Rates, error) {
	query, _ := url.Parse(c.baseURL.String())
	query.Path += "/forex/rates"
	params := url.Values{}
	params.Add("base", base)
	query.RawQuery = params.Encode()

	res, err := c.Get(query.String())
	if err != nil {
		return &Rates{}, err
	}

	rates := &Rates{}
	err = json.NewDecoder(res.Body).Decode(rates)
	return rates, err
}

type RatesTask struct {
	BaseRate string
	Executor func(base string) (*Rates, error)
}

func (r *RatesTask) Execute() (interface{}, error) {
	return r.Executor(r.BaseRate)
}

func (r *RatesTask) GetParameter() string {
	return r.BaseRate
}

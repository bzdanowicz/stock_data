package finnhub

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type Quote struct {
	CurrentPrice       float32 `json:"c"`
	LowPrice           float32 `json:"l"`
	HighPrice          float32 `json:"h"`
	OpenPrice          float32 `json:"o"`
	PreviousClosePrice float32 `json:"pc"`
}

func (c *Client) GetQuote(symbol string) (*Quote, error) {
	time.Sleep(time.Second * 1)
	query, _ := url.Parse(c.baseURL.String())
	query.Path += "/quote"
	params := url.Values{}
	params.Add("symbol", symbol)
	query.RawQuery = params.Encode()

	fmt.Println(query)
	res, err := c.Get(query.String())
	if err != nil {
		return &Quote{}, err
	}

	quote := &Quote{}
	err = json.NewDecoder(res.Body).Decode(quote)
	return quote, err
}

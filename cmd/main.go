package main

import (
	"fmt"

	"github.com/bzdanowicz/stock_data/finnhub"
	"github.com/bzdanowicz/stock_data/workerpool"
)

type QuoteTask struct {
	symbol   string
	executor func(symbol string) (*finnhub.Quote, error)
}

func (q *QuoteTask) Perform() (interface{}, error) {
	return q.executor(q.symbol)
}

func main() {
	fmt.Println("Stack Data provider.")

	c := finnhub.NewClient("api-key")

	dispatcher := workerpool.NewDispatcher(2)
	dispatcher.Start()

	task := QuoteTask{"AAPL", c.GetQuote}
	dispatcher.Enqueue(&task)
	quote := dispatcher.GetResult().TaskResult

	fmt.Println(quote)
}

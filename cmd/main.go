package main

import (
	"fmt"
	"time"

	"github.com/bzdanowicz/stock_data/finnhub"
	"github.com/bzdanowicz/stock_data/workerpool"
)

func refreshData(ticker *time.Ticker) {
	for {
		<-ticker.C
		fmt.Println("tick")
	}
}

func main() {
	fmt.Println("Stack Data provider.")

	c := finnhub.NewClient("bui1usn48v6rfhsb6kp0")

	dispatcher := workerpool.NewDispatcher(4, 50)
	dispatcher.Start()

	task := finnhub.QuoteTask{Symbol: "AAPL", Executor: c.GetQuote}
	dispatcher.Enqueue(&task)
	dispatcher.WaitAllFinished()
	quote := dispatcher.GetResult().TaskResult
	fmt.Println(quote)

	ticker := time.NewTicker(5 * time.Second)
	refreshData(ticker)
}

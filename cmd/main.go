package main

import (
	"log"
	"time"

	"github.com/bzdanowicz/stock_data/finnhub"
	"github.com/bzdanowicz/stock_data/internal"
	"github.com/bzdanowicz/stock_data/workerpool"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func updateAllTables(dispatcher *workerpool.Dispatcher, client *finnhub.Client, data *internal.Data, quoteTable *widgets.Table) {
	internal.RequestQuotes(data.Quotes, dispatcher, client)
	dispatcher.WaitAllFinished()
	internal.UpdateTable(quoteTable, data.Quotes)
	ui.Render(quoteTable)
}

func main() {
	configuration := internal.ReadConfiguration()

	data := internal.Data{}
	data.Initialize(configuration.UserQuotes)

	client := finnhub.NewClient(configuration.ApiKey)

	dispatcher := workerpool.NewDispatcher(4, 50)
	dispatcher.Start()
	go internal.DataReader(&data, dispatcher)

	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := internal.CreateTable("Stock Quotes", len(configuration.UserQuotes))

	updateAllTables(dispatcher, client, &data, table)
	ui.Render(table)

	ticker := time.NewTicker(10 * time.Second)

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker.C:
			updateAllTables(dispatcher, client, &data, table)
			ui.Render(table)
		}

	}
}

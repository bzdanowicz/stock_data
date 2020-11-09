package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bzdanowicz/stock_data/finnhub"
	"github.com/bzdanowicz/stock_data/internal"
	"github.com/bzdanowicz/stock_data/workerpool"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func updateAllTables(dispatcher *workerpool.Dispatcher, client *finnhub.Client, data *internal.Data, quoteTable *widgets.Table, ratesTable *widgets.Table) {
	internal.RequestQuotes(data.Quotes, dispatcher, client)
	internal.RequestRates(data.Rates.Base, dispatcher, client)
	dispatcher.WaitAllFinished()
	internal.UpdateTable(quoteTable, data.Quotes)
	internal.UpdateTable(ratesTable, data.Rates)
	ui.Render(quoteTable)
	ui.Render(ratesTable)
}

func handleMouseClick(mouseEvent *ui.Mouse, quotesTable *widgets.Table) {
	index, data := internal.GetRecordFromCoordinates(quotesTable, mouseEvent)
	if data == nil {
		return
	}
	row := data.([]string)
	if index == 0 {
		return
	}
	quote_symbol := row[0]
	fmt.Println(quote_symbol)
}

func main() {
	configuration := internal.ReadConfiguration()
	data := internal.Data{}
	data.Initialize(configuration.UserQuotes, configuration.BaseCurrency)

	client := finnhub.NewClient(configuration.ApiKey)

	dispatcher := workerpool.NewDispatcher(4, 50)
	dispatcher.Start()
	go internal.DataReader(&data, dispatcher)

	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	defer ui.Close()

	quotesTable := internal.CreateTable("Stock Quotes", len(configuration.UserQuotes), 0)
	ratesTable := internal.CreateTable("Rates", 2, quotesTable.GetRect().Max.Y)

	updateAllTables(dispatcher, client, &data, quotesTable, ratesTable)

	ticker := time.NewTicker(10 * time.Second)

	uiEvents := ui.PollEvents()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<MouseLeft>":
				mouseEvent := e.Payload.(ui.Mouse)
				handleMouseClick(&mouseEvent, quotesTable)
			}
		case <-ticker.C:
			updateAllTables(dispatcher, client, &data, quotesTable, ratesTable)
		}

	}
}

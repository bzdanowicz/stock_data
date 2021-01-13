package main

import (
	"log"
	"time"

	"github.com/bzdanowicz/stock_data/finnhub"
	"github.com/bzdanowicz/stock_data/internal"
	"github.com/bzdanowicz/stock_data/workerpool"
	ui "github.com/gizak/termui/v3"
)

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
	tabPane := internal.CreateTabPane(quotesTable.GetRect().Dx()/2, ratesTable.GetRect().Max.Y)
	candlePlot := internal.CreateCandlePlot(tabPane.GetRect().Max.Y)

	internal.UpdateData(dispatcher, client, &data, quotesTable, ratesTable, candlePlot, tabPane)

	ticker := time.NewTicker(10 * time.Second)
	uiEvents := ui.PollEvents()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "<Left>":
				tabPane.FocusLeft()
				internal.UpdatePlotData(dispatcher, client, &data, candlePlot, tabPane)
			case "<Right>":
				tabPane.FocusRight()
				internal.UpdatePlotData(dispatcher, client, &data, candlePlot, tabPane)
			case "q", "<C-c>":
				return
			case "<MouseLeft>":
				mouseEvent := e.Payload.(ui.Mouse)
				internal.HandleMouseClick(&mouseEvent, quotesTable, dispatcher, client, &data, candlePlot, tabPane)
			}

		case <-ticker.C:
			internal.UpdateData(dispatcher, client, &data, quotesTable, ratesTable, candlePlot, tabPane)
		}
	}
}

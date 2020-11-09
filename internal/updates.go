package internal

import (
	"log"

	"github.com/bzdanowicz/stock_data/finnhub"
	"github.com/bzdanowicz/stock_data/workerpool"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func RequestQuotes(quotes *QuoteData, dispatcher *workerpool.Dispatcher, client *finnhub.Client) {
	for key := range *quotes {
		task := finnhub.QuoteTask{Symbol: key, Executor: client.GetQuote}
		dispatcher.Enqueue(&task)
	}
}

func RequestRates(base string, dispatcher *workerpool.Dispatcher, client *finnhub.Client) {
	task := finnhub.RatesTask{BaseRate: base, Executor: client.GetRates}
	dispatcher.Enqueue(&task)

}

func DataReader(data *Data, dispatcher *workerpool.Dispatcher) {
	for {
		result := dispatcher.GetResult()
		switch res := result.TaskResult.(type) {
		case *finnhub.Quote:
			(data.Quotes)[(*result.RequestedTask).GetParameter()] = *res
		case *finnhub.Rates:
			data.Rates.Values = res.Values
		default:
			log.Fatalln("Unsupported type")
		}
	}
}

func UpdateData(dispatcher *workerpool.Dispatcher, client *finnhub.Client, data *Data, quoteTable *widgets.Table, ratesTable *widgets.Table) {
	RequestQuotes(&data.Quotes, dispatcher, client)
	RequestRates(data.Rates.Base, dispatcher, client)
	dispatcher.WaitAllFinished()
	UpdateTable(quoteTable, data.Quotes)
	UpdateTable(ratesTable, data.Rates)
	ui.Render(quoteTable)
	ui.Render(ratesTable)
}

func HandleMouseClick(mouseEvent *ui.Mouse, quotesTable *widgets.Table) {
	index, data := GetRecordFromCoordinates(quotesTable, mouseEvent)
	if data == nil {
		return
	}
	row := data.([]string)
	if index == 0 {
		return
	}
	quote_symbol := row[0]
	log.Println(quote_symbol)
}

package internal

import (
	"fmt"
	"time"

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

func RequestCandle(symbol string, from time.Time, to time.Time, dispatcher *workerpool.Dispatcher, client *finnhub.Client) {
	task := finnhub.CandleTask{
		Symbol:   symbol,
		From:     from,
		To:       to,
		Executor: client.GetCandle,
	}
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
		case *finnhub.Candle:
			data.Candle = *res
		}
	}
}

func UpdateData(dispatcher *workerpool.Dispatcher, client *finnhub.Client, data *Data, quoteTable *widgets.Table, ratesTable *widgets.Table, candlePlot *CandlePlot) {
	RequestQuotes(&data.Quotes, dispatcher, client)
	RequestRates(data.Rates.Base, dispatcher, client)
	dispatcher.WaitAllFinished()
	UpdateTable(quoteTable, data.Quotes)
	UpdateTable(ratesTable, data.Rates)
	ui.Render(quoteTable)
	ui.Render(ratesTable)

	if candlePlot.Active {
		to := time.Now()
		from := to.AddDate(0, 0, -7)
		RequestCandle(candlePlot.Symbol, from, to, dispatcher, client)
		dispatcher.WaitAllFinished()
		UpdatePlot(candlePlot.Plot, data.Candle)
		ui.Render(candlePlot.Plot)
	}
}

func HandleMouseClick(mouseEvent *ui.Mouse, quotesTable *widgets.Table, candlePlot *CandlePlot) {
	index, data := GetRecordFromCoordinates(quotesTable, mouseEvent)
	if data == nil {
		return
	}
	row := data.([]string)
	if index == 0 {
		return
	}
	quote_symbol := row[0]
	candlePlot.Active = true
	candlePlot.Symbol = quote_symbol
}

type CandlePlot struct {
	Plot   *SimplePlot
	Symbol string
	From   time.Time
	To     time.Time
	Active bool
}

func CreateCandlePlot(y int) *CandlePlot {
	candlePlot := CandlePlot{Plot: NewSimplePlot(), Active: false}
	candlePlot.Plot.SetRect(0, y, 150, y+30)
	return &candlePlot
}

func UpdatePlot(plot *SimplePlot, data interface{}) {
	plot.DataLabels = make([]string, 0)

	for _, t := range data.(finnhub.Candle).Timestamps {
		plot.DataLabels = append(plot.DataLabels, fmt.Sprintf("%v", t))
	}

	plot.Data = data.(finnhub.Candle).CurrentPrice

	plot.HorizontalScale = plot.Dx() / (len(plot.Data))
	if plot.HorizontalScale < 1 {
		plot.HorizontalScale = 1
	}
}

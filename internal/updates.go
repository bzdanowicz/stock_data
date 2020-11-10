package internal

import (
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

func UpdatePlotData(dispatcher *workerpool.Dispatcher, client *finnhub.Client, data *Data, candlePlot *CandlePlot) {
	if candlePlot.Active {
		to := time.Now()
		from := to.AddDate(-3, 0, 0)
		RequestCandle(candlePlot.Symbol, from, to, dispatcher, client)
		dispatcher.WaitAllFinished()
		UpdatePlot(candlePlot.Plot, data.Candle)
		if len(candlePlot.Plot.Data) == len(candlePlot.Plot.DataLabels) && len(candlePlot.Plot.Data) != 0 {
			ui.Render(candlePlot.Plot)
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

	UpdatePlotData(dispatcher, client, data, candlePlot)
}

func HandleMouseClick(mouseEvent *ui.Mouse, quotesTable *widgets.Table, dispatcher *workerpool.Dispatcher, client *finnhub.Client, data *Data, candlePlot *CandlePlot) {
	index, record := GetRecordFromCoordinates(quotesTable, mouseEvent)
	if record == nil {
		return
	}
	row := record.([]string)
	if index == 0 {
		return
	}
	quote_symbol := row[0]
	candlePlot.Active = true
	candlePlot.Symbol = quote_symbol
	UpdatePlotData(dispatcher, client, data, candlePlot)
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
	candlePlot.Plot.SetRect(0, y, 210, y+30)
	return &candlePlot
}

func UpdatePlot(plot *SimplePlot, data interface{}) {
	plot.DataLabels = make([]string, 0)
	for _, t := range data.(finnhub.Candle).Timestamps {
		tm := time.Unix(t, 0)
		plot.DataLabels = append(plot.DataLabels, tm.Format("02/Jan/06"))
	}
	plot.Data = data.(finnhub.Candle).CurrentPrice
}

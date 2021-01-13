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

func RequestCandle(symbol string, from time.Time, to time.Time, resolution string, dispatcher *workerpool.Dispatcher, client *finnhub.Client) {
	task := finnhub.CandleTask{
		Symbol:     symbol,
		From:       from,
		To:         to,
		Resolution: resolution,
		Executor:   client.GetCandle,
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

func UpdatePlotData(dispatcher *workerpool.Dispatcher, client *finnhub.Client, data *Data, candlePlot *CandlePlot, tabPane *widgets.TabPane) {
	if candlePlot.Active {
		to := time.Now()
		var from time.Time
		var resolution string
		switch tabPane.ActiveTabIndex {
		case 0:
			from = to.AddDate(0, 0, -7)
			resolution = "5"
		case 1:
			from = to.AddDate(0, -1, 0)
			resolution = "5"
		case 2:
			from = to.AddDate(-1, 0, 0)
			resolution = "D"
		case 3:
			from = to.AddDate(-5, 0, 0)
			resolution = "D"
		}

		RequestCandle(candlePlot.Symbol, from, to, resolution, dispatcher, client)
		dispatcher.WaitAllFinished()
		UpdatePlot(candlePlot.Plot, data.Candle)
		if len(candlePlot.Plot.Data) == len(candlePlot.Plot.DataLabels) && len(candlePlot.Plot.Data) != 0 {
			ui.Render(candlePlot.Plot)
			ui.Render(tabPane)
		}
	}
}
func UpdateData(dispatcher *workerpool.Dispatcher, client *finnhub.Client, data *Data, quoteTable *widgets.Table, ratesTable *widgets.Table, candlePlot *CandlePlot, tabPane *widgets.TabPane) {
	RequestQuotes(&data.Quotes, dispatcher, client)
	RequestRates(data.Rates.Base, dispatcher, client)
	dispatcher.WaitAllFinished()
	UpdateTable(quoteTable, data.Quotes)
	UpdateTable(ratesTable, data.Rates)
	ui.Render(quoteTable)
	ui.Render(ratesTable)

	UpdatePlotData(dispatcher, client, data, candlePlot, tabPane)
}

func HandleMouseClick(mouseEvent *ui.Mouse, quotesTable *widgets.Table, dispatcher *workerpool.Dispatcher, client *finnhub.Client, data *Data, candlePlot *CandlePlot, tabPane *widgets.TabPane) {
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
	candlePlot.Plot.Block.Title = quote_symbol
	UpdatePlotData(dispatcher, client, data, candlePlot, tabPane)
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

func CreateTabPane(x int, y int) *widgets.TabPane {
	width := 32
	tabpane := widgets.NewTabPane("Week", "Month", "Year", "5 Years")
	tabpane.SetRect(x-width/2, y, x+width/2-1, y+3)
	tabpane.Border = true
	return tabpane
}

package internal

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/bzdanowicz/stock_data/finnhub"
	"github.com/bzdanowicz/stock_data/workerpool"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type QuoteData map[string]finnhub.Quote

type Data struct {
	Quotes QuoteData
}

func (data *Data) Initialize(quotes []string) {
	data.Quotes = make(QuoteData)
	for _, q := range quotes {
		data.Quotes[q] = finnhub.Quote{}
	}
}

func CreateTable(title string, rows int) *widgets.Table {
	table := widgets.NewTable()
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.RowSeparator = true
	table.TextAlignment = ui.AlignCenter
	table.Block.Title = title
	table.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
	table.SetRect(0, 0, 210, 3+rows*2)
	return table
}

func UpdateTable(table *widgets.Table, data interface{}) {
	switch data.(type) {
	case QuoteData:
		table.Rows = [][]string{
			[]string{"Symbol", "Current price", "Open price of the day", "Low price of the day", "High price of the day", "Previous close price", "Change %"},
		}
		for symbol, quote := range data.(QuoteData) {
			change := (quote.CurrentPrice - quote.PreviousClosePrice) / quote.PreviousClosePrice * 100
			table.Rows = append(table.Rows, []string{symbol, fmt.Sprintf("%f", quote.CurrentPrice), fmt.Sprintf("%f", quote.OpenPrice),
				fmt.Sprintf("%f", quote.LowPrice), fmt.Sprintf("%f", quote.HighPrice), fmt.Sprintf("%f", quote.PreviousClosePrice), fmt.Sprintf("%f", change)})

		}
		sort.SliceStable(table.Rows[1:], func(i, j int) bool {
			return table.Rows[i+1][0] < table.Rows[j+1][0]
		})

		for i := range table.Rows {
			val, err := strconv.ParseFloat(table.Rows[i][6], 32)
			if err != nil {
				continue
			}
			if val > 0 {
				table.RowStyles[i] = ui.NewStyle(ui.ColorGreen, ui.ColorBlack, ui.ModifierBold)
			} else {
				table.RowStyles[i] = ui.NewStyle(ui.ColorRed, ui.ColorBlack, ui.ModifierBold)
			}
		}
	}
}

func RequestQuotes(quotes QuoteData, dispatcher *workerpool.Dispatcher, client *finnhub.Client) {
	for key := range quotes {
		task := finnhub.QuoteTask{Symbol: key, Executor: client.GetQuote}
		dispatcher.Enqueue(&task)
	}
}

func DataReader(data *Data, dispatcher *workerpool.Dispatcher) {
	for {
		result := dispatcher.GetResult()
		switch res := result.TaskResult.(type) {
		case *finnhub.Quote:
			(data.Quotes)[(*result.RequestedTask).GetParameter()] = *res
		default:
			log.Fatalln("Unsupported type")
		}
	}
}

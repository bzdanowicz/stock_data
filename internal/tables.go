package internal

import (
	"fmt"
	"image"
	"sort"
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	tableWidth = 210
)

func CreateTable(title string, rows int, y int) *widgets.Table {
	table := widgets.NewTable()
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.RowSeparator = true
	table.TextAlignment = ui.AlignCenter
	table.Block.Title = title
	table.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
	rows += 1 // We also need to include Description row
	table.SetRect(0, y, tableWidth, y+rows*2+1)
	return table
}

func GetRecordFromCoordinates(table *widgets.Table, mouseEvent *ui.Mouse) (int, interface{}) {
	overlaps := table.GetRect().Overlaps(image.Rect(mouseEvent.X, mouseEvent.Y, mouseEvent.X+1, mouseEvent.Y+1))
	if !overlaps {
		return 0, nil
	}

	verticalPoint := mouseEvent.Y - table.GetRect().Min.Y
	calculatedIndex := (verticalPoint - 1) / 2
	return calculatedIndex, table.Rows[calculatedIndex]
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
	case RatesData:
		rates := data.(RatesData)
		ratesValues := rates.Values
		table.Rows = [][]string{
			[]string{"Symbol", "USD", "EUR", "GBP", "CHF", "JPY", "PLN"},
			[]string{rates.Base, fmt.Sprintf("%f", ratesValues.USD), fmt.Sprintf("%f", ratesValues.EUR), fmt.Sprintf("%f", ratesValues.GBP),
				fmt.Sprintf("%f", ratesValues.CHF), fmt.Sprintf("%f", ratesValues.JPY), fmt.Sprintf("%f", ratesValues.PLN)},
			[]string{fmt.Sprintf("Inverse currency rate %s", rates.Base), fmt.Sprintf("%f", 1/ratesValues.USD), fmt.Sprintf("%f", 1/ratesValues.EUR),
				fmt.Sprintf("%f", 1/ratesValues.GBP), fmt.Sprintf("%f", 1/ratesValues.CHF), fmt.Sprintf("%f", 1/ratesValues.JPY), fmt.Sprintf("%f", 1/ratesValues.PLN)},
		}
	}
}

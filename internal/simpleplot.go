package internal

import (
	"fmt"
	"image"

	"github.com/gizak/termui/v3"
	. "github.com/gizak/termui/v3"
)

type SimplePlot struct {
	Block

	Data       []float64
	DataLabels []string
	MaxVal     float64
	MinVal     float64

	LineColors []Color
	AxesColor  Color
	ShowAxes   bool

	HorizontalScale int
}

const (
	xAxisLabelsHeight = 1
	yAxisLabelsWidth  = 2
	xAxisLabelsGap    = 0
	yAxisLabelsGap    = 1
)

func NewSimplePlot() *SimplePlot {
	return &SimplePlot{
		Block:           *NewBlock(),
		LineColors:      Theme.Plot.Lines,
		AxesColor:       Theme.Plot.Axes,
		Data:            []float64{},
		HorizontalScale: 1,
		ShowAxes:        true,
	}
}

/*func (self *SimplePlot) renderBraille(buf *Buffer, drawArea image.Rectangle, maxVal float64) {
	canvas := NewCanvas()
	canvas.Rectangle = drawArea

	for i, line := range self.Data {
		previousHeight := int((line[1] / maxVal) * float64(drawArea.Dy()-1))
		for j, val := range line[1:] {
			height := int((val / maxVal) * float64(drawArea.Dy()-1))
			canvas.SetLine(
				image.Pt(
					(drawArea.Min.X+(j*self.HorizontalScale))*2,
					(drawArea.Max.Y-previousHeight-1)*4,
				),
				image.Pt(
					(drawArea.Min.X+((j+1)*self.HorizontalScale))*2,
					(drawArea.Max.Y-height-1)*4,
				),
				SelectColor(self.LineColors, i),
			)
			previousHeight = height
		}
	}

	canvas.Draw(buf)
}*/

func (self *SimplePlot) renderDot(buf *Buffer, drawArea image.Rectangle, maxVal float64) {
	for j, line := range self.Data {
		val := line
		height := int((val / maxVal) * float64(drawArea.Dy()-1))
		buf.SetCell(
			NewCell(termui.DOT, NewStyle(SelectColor(self.LineColors, j))),
			image.Pt(drawArea.Min.X+(j*self.HorizontalScale), drawArea.Max.Y-1-height),
		)
	}
}

func (self *SimplePlot) plotAxes(buf *Buffer, maxVal float64) {
	// draw origin cell
	buf.SetCell(
		NewCell(BOTTOM_LEFT, NewStyle(ColorWhite)),
		image.Pt(self.Inner.Min.X+yAxisLabelsWidth, self.Inner.Max.Y-xAxisLabelsHeight-1),
	)
	// draw x axis line
	for i := yAxisLabelsWidth + 1; i < self.Inner.Dx(); i++ {
		buf.SetCell(
			NewCell(HORIZONTAL_DASH, NewStyle(ColorWhite)),
			image.Pt(i+self.Inner.Min.X, self.Inner.Max.Y-xAxisLabelsHeight-1),
		)
	}
	// draw y axis line
	for i := 0; i < self.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buf.SetCell(
			NewCell(VERTICAL_DASH, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Min.X+yAxisLabelsWidth, i+self.Inner.Min.Y),
		)
	}
	// draw x axis labels
	// draw 0
	buf.SetString(
		self.DataLabels[0],
		NewStyle(ColorWhite),
		image.Pt(self.Inner.Min.X+yAxisLabelsWidth, self.Inner.Max.Y-1),
	)
	// draw rest
	for x := self.Inner.Min.X + yAxisLabelsWidth + (xAxisLabelsGap)*self.HorizontalScale + 1; x < self.Inner.Max.X-1; {
		label_name := self.DataLabels[(x-(self.Inner.Min.X+yAxisLabelsWidth)-1)/(self.HorizontalScale)+1]
		label := fmt.Sprintf(
			"%d",
			(x-(self.Inner.Min.X+yAxisLabelsWidth)-1)/(self.HorizontalScale)+1,
		)
		buf.SetString(
			label_name,
			NewStyle(ColorWhite),
			image.Pt(x, self.Inner.Max.Y-1),
		)
		x += (len(label) + xAxisLabelsGap) * self.HorizontalScale
	}
	// draw y axis labels
	verticalScale := maxVal / float64(self.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 0; i*(yAxisLabelsGap+1) < self.Inner.Dy()-1; i++ {
		buf.SetString(
			fmt.Sprintf("%.2f", float64(i)*verticalScale*(yAxisLabelsGap+1)),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Min.X, self.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2),
		)
	}
}

func (self *SimplePlot) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	maxVal := self.MaxVal
	minVal := self.MinVal
	if maxVal == 0 || minVal == 0 {
		for _, val := range self.Data {
			if val > maxVal {
				maxVal = val
			}
			if val < minVal {
				minVal = val
			}
		}
	}

	if self.ShowAxes {
		self.plotAxes(buf, maxVal)
	}

	drawArea := self.Inner
	if self.ShowAxes {
		drawArea = image.Rect(
			self.Inner.Min.X+yAxisLabelsWidth+1, self.Inner.Min.Y,
			self.Inner.Max.X, self.Inner.Max.Y-xAxisLabelsHeight-1,
		)
	}

	self.renderDot(buf, drawArea, maxVal)
}

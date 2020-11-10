package internal

import (
	"fmt"
	"image"

	"github.com/gizak/termui/v3"
	. "github.com/gizak/termui/v3"
)

type PlotMarker uint

const (
	MarkerBraille PlotMarker = iota
	MarkerDot
)

type SimplePlot struct {
	Block

	Data       []float64
	DataLabels []string

	Marker     PlotMarker
	LineColors []Color
	AxesColor  Color
	ShowAxes   bool

	horizontalScale float32
}

const (
	xAxisLabelsHeight = 1
	yAxisLabelsWidth  = 2
	xAxisLabelsGap    = 1
	yAxisLabelsGap    = 1
)

func NewSimplePlot() *SimplePlot {
	return &SimplePlot{
		Block:      *NewBlock(),
		LineColors: Theme.Plot.Lines,
		AxesColor:  Theme.Plot.Axes,
		Data:       []float64{},
		ShowAxes:   true,
		Marker:     MarkerDot,
	}
}

func (self *SimplePlot) renderDot(buf *Buffer, drawArea image.Rectangle, minVal float64, maxVal float64) {
	h := maxVal - minVal

	for j, line := range self.Data {
		val := line
		height := int(((val - minVal) / h) * float64(drawArea.Dy()-1))
		buf.SetCell(
			NewCell(termui.DOT, NewStyle(ColorGreen)),
			image.Pt(drawArea.Min.X+int(float32(j)*self.horizontalScale), drawArea.Max.Y-1-height),
		)
	}
}

func (self *SimplePlot) renderBraille(buf *Buffer, drawArea image.Rectangle, minVal float64, maxVal float64) {
	canvas := NewCanvas()
	canvas.Rectangle = drawArea
	h := maxVal - minVal
	prevPoint := image.Point{}

	for j, line := range self.Data {
		val := line
		height := int(((val - minVal) / h) * float64(drawArea.Dy()-1))
		point := image.Pt(drawArea.Min.X+int(float32(j)*self.horizontalScale)*2, (drawArea.Max.Y-1-height)*4)
		canvas.SetLine(point, prevPoint, ColorWhite)
		prevPoint = point
	}
	canvas.Draw(buf)
}

func (self *SimplePlot) plotAxes(buf *Buffer, minVal float64, maxVal float64) {
	buf.SetCell(
		NewCell(BOTTOM_LEFT, NewStyle(ColorWhite)),
		image.Pt(self.Inner.Min.X+yAxisLabelsWidth, self.Inner.Max.Y-xAxisLabelsHeight-1),
	)

	for i := yAxisLabelsWidth + 1; i < self.Inner.Dx(); i++ {
		buf.SetCell(
			NewCell(HORIZONTAL_DASH, NewStyle(ColorWhite)),
			image.Pt(i+self.Inner.Min.X, self.Inner.Max.Y-xAxisLabelsHeight-1),
		)
	}

	for i := 0; i < self.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buf.SetCell(
			NewCell(VERTICAL_DASH, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Min.X+yAxisLabelsWidth, i+self.Inner.Min.Y),
		)
	}

	for x := self.Inner.Min.X + yAxisLabelsWidth + int((float32(xAxisLabelsGap))*self.horizontalScale) + 1; x < self.Inner.Max.X-5; {
		index := len(self.DataLabels) * x / self.Inner.Max.X
		label_name := self.DataLabels[index]
		buf.SetString(
			label_name,
			NewStyle(ColorWhite),
			image.Pt(x, self.Inner.Max.Y-1),
		)
		x += len(label_name) + xAxisLabelsGap
	}

	h := maxVal - minVal
	verticalScale := h / float64(self.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 0; i*(yAxisLabelsGap+1) < self.Inner.Dy()-1; i++ {
		yVal := minVal + float64(i)*verticalScale*(yAxisLabelsGap+1)
		buf.SetString(
			fmt.Sprintf("%.2f", yVal),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Min.X, self.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2),
		)
	}
}

func (self *SimplePlot) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	maxVal := float64(0)
	minVal := self.Data[0]

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

	self.horizontalScale = float32(self.Inner.Dx()-5) / float32(len(self.Data))

	if self.ShowAxes {
		self.plotAxes(buf, minVal, maxVal)
	}

	drawArea := self.Inner
	if self.ShowAxes {
		drawArea = image.Rect(
			self.Inner.Min.X+yAxisLabelsWidth+1, self.Inner.Min.Y,
			self.Inner.Max.X, self.Inner.Max.Y-xAxisLabelsHeight-1,
		)
	}

	self.renderDot(buf, drawArea, minVal, maxVal)
}

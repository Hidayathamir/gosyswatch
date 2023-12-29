package main

import (
	"flag"

	"github.com/Hidayathamir/gosyswatch/chart"
	"github.com/Hidayathamir/gosyswatch/proc"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

func main() {
	var lineChartWidth int
	flag.IntVar(&lineChartWidth, "width", 50, "line chart width")
	flag.Parse()

	app := tview.NewApplication()

	cpusUsage, err := proc.GetCPUsUsageInPercentage()
	if err != nil {
		panic(err)
	}

	cpuPlotChart := createPlotChart("CPU")
	cpuPlotChart.SetLineColor(generateCPULineColors(cpusUsage))

	memoryPlotChart := createPlotChart("Memory")

	go func() {
		cpuLineCharts := createLineCharts(lineChartWidth, len(cpusUsage))

		memoryChart := chart.NewSizedLineChart(lineChartWidth)
		memoryUsageHistory := [][]float64{memoryChart.Values}

		for {
			updateCPUCharts(cpuLineCharts)
			cpuPlotChart.SetData(collectChartData(cpuLineCharts))

			updateMemoryChart(&memoryChart, memoryUsageHistory)
			memoryPlotChart.SetData(memoryUsageHistory)

			app.Draw()
		}
	}()

	layout := createLayout(cpuPlotChart, memoryPlotChart)
	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// generateCPULineColors generates line colors based on CPU number.
func generateCPULineColors(cpusUsage []float64) []tcell.Color {
	lineColors := []tcell.Color{}
	for i := range cpusUsage {
		lineColors = append(lineColors, tcell.ColorMaroon+tcell.Color(i))
	}
	return lineColors
}

// createPlotChart creates a new TVX widget plot chart with the given title.
func createPlotChart(title string) *tvxwidgets.Plot {
	plotChart := tvxwidgets.NewPlot()
	plotChart.SetBorder(true)
	plotChart.SetTitle(title)
	plotChart.SetMarker(tvxwidgets.PlotMarkerBraille)
	return plotChart
}

// collectChartData collects chart data from SizedLineChart instances
func collectChartData(charts []chart.SizedLineChart) [][]float64 {
	data := make([][]float64, len(charts))

	for i := 0; i < len(charts); i++ {
		chartData := charts[i]
		data[i] = chartData.Values
	}

	return data
}

// createLineCharts creates multiple SizedLineChart instances with the given width and count.
func createLineCharts(width, count int) []chart.SizedLineChart {
	lineCharts := make([]chart.SizedLineChart, count)
	for i := range lineCharts {
		lineCharts[i] = chart.NewSizedLineChart(width)
	}
	return lineCharts
}

// updateCPUCharts updates the CPU line charts with the current CPU usage data.
func updateCPUCharts(cpuLineCharts []chart.SizedLineChart) {
	cpusUsage, err := proc.GetCPUsUsageInPercentage()
	if err != nil {
		panic(err)
	}

	for i, cpuUsage := range cpusUsage {
		cpuLineCharts[i].Add(cpuUsage)
	}
}

// updateMemoryChart updates the memory line chart with the current memory usage data.
func updateMemoryChart(memoryChart *chart.SizedLineChart, memoryUsageHistory [][]float64) {
	memUsage, err := proc.GetMemoryUsagePercentage()
	if err != nil {
		panic(err)
	}
	memoryChart.Add(memUsage)
	memoryUsageHistory[0] = memoryChart.Values
}

// createLayout creates a layout with CPU and Memory plot charts.
func createLayout(cpuPlotChart, memoryPlotChart *tvxwidgets.Plot) *tview.Flex {
	firstRow := tview.NewFlex().AddItem(cpuPlotChart, 0, 1, false)
	secondRow := tview.NewFlex().AddItem(memoryPlotChart, 0, 1, false)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(firstRow, 0, 1, false).
		AddItem(secondRow, 0, 1, false)
	return layout
}

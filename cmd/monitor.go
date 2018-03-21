package cmd

// import (
// 	"math"

// 	"github.com/gizak/termui"
// 	"github.com/spf13/cobra"
// )

// // Monitor is the dashboard to monitor everything
// var Monitor = &cobra.Command{
// 	Use:               "monitor",
// 	Short:             "Monitor the running application",
// 	DisableAutoGenTag: true,
// 	Run:               monitorHandler,
// }

// func monitorHandler(cmd *cobra.Command, args []string) {
// 	if err := termui.Init(); err != nil {
// 		panic(err)
// 	}
// 	defer termui.Close()

// 	p := termui.NewPar(":PRESS q TO QUIT DEMO")
// 	p.Height = 3
// 	p.Width = 50
// 	p.BorderLabel = "Text Box"

// 	strs := []string{"[0] gizak/termui", "[1] editbox.go", "[2] interrupt.go", "[3] keyboard.go", "[4] output.go", "[5] random_out.go", "[6] dashboard.go", "[7] nsf/termbox-go"}
// 	list := termui.NewList()
// 	list.Items = strs
// 	list.BorderLabel = "List"
// 	list.Height = 10

// 	g := termui.NewGauge()
// 	g.Percent = 50
// 	g.Height = 3
// 	g.BorderLabel = "Gauge"

// 	spark := termui.Sparkline{}
// 	spark.Height = 1
// 	spark.Title = "srv 0:"
// 	spdata := []int{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}
// 	spark.Data = spdata

// 	spark1 := termui.Sparkline{}
// 	spark1.Height = 1
// 	spark1.Title = "srv 1:"
// 	spark1.Data = spdata

// 	sp := termui.NewSparklines(spark, spark1)
// 	sp.Height = 10
// 	sp.BorderLabel = "Sparkline"

// 	sinps := (func() []float64 {
// 		n := 220
// 		ps := make([]float64, n)
// 		for i := range ps {
// 			ps[i] = 1 + math.Sin(float64(i)/5)
// 		}
// 		return ps
// 	})()

// 	lc := termui.NewLineChart()
// 	lc.BorderLabel = "dot-mode Line Chart"
// 	lc.Data = sinps
// 	lc.Height = 20
// 	lc.Mode = "dot"

// 	bc := termui.NewBarChart()
// 	bcdata := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
// 	bc.BorderLabel = "Bar Chart"
// 	bc.Height = 10
// 	bc.DataLabels = []string{"S0", "S1", "S2", "S3", "S4", "S5"}

// 	lc1 := termui.NewLineChart()
// 	lc1.BorderLabel = "braille-mode Line Chart"
// 	lc1.Data = sinps
// 	lc1.Height = 20

// 	p1 := termui.NewPar("Hey!\nI am a borderless block!")
// 	p1.Border = false

// 	termui.Body.AddRows(
// 		termui.NewRow(
// 			termui.NewCol(4, 0, list),
// 			termui.NewCol(4, 0, sp),
// 			termui.NewCol(4, 0, bc)),
// 		termui.NewRow(
// 			termui.NewCol(8, 0, g),
// 			termui.NewCol(4, 0, p1)),
// 		termui.NewRow(
// 			termui.NewCol(8, 0, lc),
// 			termui.NewCol(4, 0, lc1)),
// 		termui.NewRow(
// 			termui.NewCol(12, 0, p)))

// 	// calculate layout
// 	termui.Body.Align()

// 	draw := func(t int) {
// 		g.Percent = t % 101
// 		list.Items = strs[t%9:]
// 		sp.Lines[0].Data = spdata[:30+t%50]
// 		sp.Lines[1].Data = spdata[:35+t%50]
// 		lc.Data = sinps[t/2%220:]
// 		lc1.Data = sinps[2*t%220:]
// 		bc.Data = bcdata[t/2%10:]
// 		termui.Render(p, list, g, sp, lc, bc, lc1, p1)
// 	}
// 	termui.Handle("/sys/kbd/C-x", func(termui.Event) {
// 		termui.StopLoop()
// 	})
// 	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
// 		termui.StopLoop()
// 	})
// 	termui.Handle("/sys/kbd/q", func(termui.Event) {
// 		termui.StopLoop()
// 	})
// 	termui.Handle("/timer/1s", func(e termui.Event) {
// 		t := e.Data.(termui.EvtTimer)
// 		draw(int(t.Count))
// 	})
// 	termui.Loop()
// }

// func init() {
// 	RootCmd.AddCommand(Monitor)
// }

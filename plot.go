// plot.go
package main

import (
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// ---------------------------------------------------------------------
//
//	Linear regression result – used by the scatter‑+‑line plot
//
// ---------------------------------------------------------------------
type LinearFitResult struct {
	Slope     float64 // Mt CO₂ per year
	Intercept float64
	R2        float64
}

// ---------------------------------------------------------------------
//
//	PlotTimeSeries – scatter of (Year, CO₂) + regression line
//
// ---------------------------------------------------------------------
func PlotTimeSeries(records []CO2Record, fit LinearFitResult, outPath string) error {
	// 0️⃣ New plot (Gonum 0.16 returns only *Plot)
	p := plot.New()
	p.Title.Text = "Germany CO₂ Emissions (1970‑2022)"
	p.X.Label.Text = "Year"
	p.Y.Label.Text = "CO₂ (Mt)"

	// 1️⃣ Scatter points
	pts := make(plotter.XYs, len(records))
	for i, r := range records {
		pts[i].X = float64(r.Year)
		pts[i].Y = r.CO2Emissions
	}
	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		return err
	}
	scatter.GlyphStyle.Color = plotutil.Color(0) // first colour in palette
	scatter.GlyphStyle.Radius = vg.Points(2)
	p.Add(scatter)

	// 2️⃣ Regression line (straight line through first & last point)
	linePts := plotter.XYs{
		{X: pts[0].X, Y: fit.Slope*pts[0].X + fit.Intercept},
		{X: pts[len(pts)-1].X, Y: fit.Slope*pts[len(pts)-1].X + fit.Intercept},
	}
	line, err := plotter.NewLine(linePts)
	if err != nil {
		return err
	}
	line.Color = plotutil.Color(1) // second colour
	line.Width = vg.Points(1.5)
	p.Add(line)

	// 3️⃣ Legend (no dummy entries – avoids the nil‑pointer crash)
	p.Legend.Add("Data points", scatter)
	p.Legend.Add("Linear fit", line)

	// 4️⃣ Save PNG
	if err := p.Save(8*vg.Inch, 5*vg.Inch, outPath); err != nil {
		return err
	}
	log.Printf("✅ Scatter‑line plot saved to %s", outPath)
	return nil
}

// ---------------------------------------------------------------------
//
//	PlotHistogram – distribution of the CO₂ emission values
//
// ---------------------------------------------------------------------
// outPath is the filename you want (e.g. "co2_histogram.png")
func PlotHistogram(records []CO2Record, outPath string) error {
	// 1️⃣ Extract emission values into a slice that plotter expects.
	vals := make(plotter.Values, len(records))
	for i, r := range records {
		vals[i] = r.CO2Emissions
	}

	// 2️⃣ New plot for the histogram
	p := plot.New()
	p.Title.Text = "Histogram of Germany CO₂ Emissions (1970‑2022)"
	p.X.Label.Text = "CO₂ Emissions (Mt)"
	p.Y.Label.Text = "Frequency"

	// 3️⃣ Build histogram – 20 bins is a good default.
	//    Change the second argument to any positive integer if you want
	//    more or fewer bins.
	hist, err := plotter.NewHist(vals, 20)
	if err != nil {
		return err
	}
	hist.FillColor = plotutil.Color(2) // third colour in Gonum palette
	hist.LineStyle.Width = vg.Points(0.5)

	// 4️⃣ Add histogram to the plot
	p.Add(hist)

	// 5️⃣ Save PNG
	if err := p.Save(8*vg.Inch, 5*vg.Inch, outPath); err != nil {
		return err
	}
	log.Printf("✅ Histogram saved to %s", outPath)
	return nil
}

// main.go
package main

import (
	"log"
	"os"
	"path/filepath"
)

// ---------------------------------------------------------------
// Simple linear regression (pure Go – works with any Go version)
// ---------------------------------------------------------------
func linearRegression(x, y []float64) (slope, intercept, r2 float64) {
	n := float64(len(x))
	if n == 0 {
		return 0, 0, 0
	}
	var sumX, sumY, sumXY, sumXX, sumYY float64
	for i := 0; i < int(n); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumXX += x[i] * x[i]
		sumYY += y[i] * y[i]
	}
	// slope & intercept
	slope = (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)
	intercept = (sumY - slope*sumX) / n

	// R² (coefficient of determination)
	var ssRes, ssTot float64
	meanY := sumY / n
	for i := 0; i < int(n); i++ {
		yPred := slope*x[i] + intercept
		ssRes += (y[i] - yPred) * (y[i] - yPred)
		ssTot += (y[i] - meanY) * (y[i] - meanY)
	}
	if ssTot == 0 {
		r2 = 1 // perfect fit (degenerate case)
	} else {
		r2 = 1 - ssRes/ssTot
	}
	return
}

// ---------------------------------------------------------------
// Helper – safe preview of the first N rows (avoids panic)
// ---------------------------------------------------------------
func printFirstN(records []CO2Record, n int) {
	if len(records) == 0 {
		log.Println("⚠️ No records – the CSV is empty or all rows were skipped.")
		return
	}
	if n > len(records) {
		n = len(records)
	}
	log.Printf("First %d record(s):", n)
	for i := 0; i < n; i++ {
		r := records[i]
		log.Printf("[%d] %s (%s) | %d | %.2f Mt CO₂e",
			i+1, r.CountryName, r.CountryCode, r.Year, r.CO2Emissions)
	}
}

// ---------------------------------------------------------------
// Main program
// ---------------------------------------------------------------
func main() {
	// -----------------------------------------------------------------
	// 1️⃣ Load the already‑cleaned CSV
	// -----------------------------------------------------------------
	cleanedPath := "data/cleaned/germany_co2_cleaned.csv"
	cleanedData, err := LoadAndClean(cleanedPath)
	if err != nil {
		log.Fatalf("❌ Failed to load cleaned data: %v", err)
	}
	log.Printf("✅ Loaded %d cleaned records from %s", len(cleanedData), cleanedPath)

	// -----------------------------------------------------------------
	// 2️⃣ Quick sanity‑check – print first three rows (safe)
	// -----------------------------------------------------------------
	printFirstN(cleanedData, 3)

	// -----------------------------------------------------------------
	// If there are no rows, stop early – nothing to plot.
	// -----------------------------------------------------------------
	if len(cleanedData) == 0 {
		log.Println("⏹️ Exiting because the dataset is empty.")
		return
	}

	// -----------------------------------------------------------------
	// 3️⃣ Prepare slices for regression (Year → CO₂)
	// -----------------------------------------------------------------
	years := make([]float64, len(cleanedData))
	co2 := make([]float64, len(cleanedData))
	for i, r := range cleanedData {
		years[i] = float64(r.Year)
		co2[i] = r.CO2Emissions
	}

	// -----------------------------------------------------------------
	// 4️⃣ Compute slope, intercept, and R²
	// -----------------------------------------------------------------
	slope, intercept, r2 := linearRegression(years, co2)
	log.Printf("📈 Regression → slope: %.3f Mt/yr, intercept: %.3f, R²: %.4f",
		slope, intercept, r2)

	// -----------------------------------------------------------------
	// 5️⃣ Package the regression numbers for the plot helper
	// -----------------------------------------------------------------
	fit := LinearFitResult{
		Slope:     slope,
		Intercept: intercept,
		R2:        r2,
	}

	// -----------------------------------------------------------------
	// 6️⃣ Draw the scatter + regression line
	// -----------------------------------------------------------------
	if err := PlotTimeSeries(cleanedData, fit, "co2_trend.png"); err != nil {
		log.Fatalf("❌ Plot generation failed: %v", err)
	}
	log.Println("✅ Plot saved to co2_trend.png – open it in the VS Code Explorer pane")

	// -----------------------------------------------------------------
	// 7️⃣ **NEW** – draw a histogram of the CO₂ emission values
	// -----------------------------------------------------------------
	if err := PlotHistogram(cleanedData, "co2_histogram.png"); err != nil {
		log.Fatalf("❌ Histogram generation failed: %v", err)
	}
	log.Println("✅ Histogram saved to co2_histogram.png – open it in VS Code Explorer")

	// -----------------------------------------------------------------
	// 8️⃣ (Optional) write a copy of the cleaned CSV back out
	// -----------------------------------------------------------------
	outputPath := "data/cleaned/germany_co2_cleaned_copy.csv"
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		log.Fatalf("❌ Failed to create output directory: %v", err)
	}
	if err := SaveCleaned(cleanedData, outputPath); err != nil {
		log.Fatalf("❌ Failed to write cleaned CSV copy: %v", err)
	}
	log.Printf("✅ Cleaned CSV copy written to %s", outputPath)
}

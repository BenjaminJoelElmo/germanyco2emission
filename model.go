package main

import (
	"log"
	"math"
)

func AnalyzeTrend(records []CO2Record) {
	years := make([]float64, len(records))
	co2 := make([]float64, len(records))

	for i, r := range records {
		years[i] = float64(r.Year)
		co2[i] = r.CO2Emissions
	}

	// Manual regression (no external dependencies)
	n := float64(len(years))
	var sumX, sumY, sumXY, sumX2 float64
	for i := range years {
		x := years[i]
		y := co2[i]
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	denominator := n*sumX2 - sumX*sumX
	if denominator == 0 {
		log.Fatal("Cannot compute regression: all years are identical")
	}

	slope := (n*sumXY - sumX*sumY) / denominator
	intercept := (sumY - slope*sumX) / n

	// R² calculation
	var ssTotal, ssResidual float64
	co2Mean := sumY / n
	for i := range co2 {
		yPred := slope*years[i] + intercept
		ssTotal += math.Pow(co2[i]-co2Mean, 2)
		ssResidual += math.Pow(co2[i]-yPred, 2)
	}
	r2 := 1 - (ssResidual / ssTotal)
	if r2 < 0 {
		r2 = 0
	}

	// Climate interpretation
	log.Printf("📈 CO₂ Trend Analysis (1970-2022):")
	log.Printf("• Slope: %.2f Mt CO₂/year", slope)
	log.Printf("• R²: %.2f", r2)

	if slope < 0 {
		log.Printf("• ✅ Germany reduced emissions by %.1f Mt/year", math.Abs(slope))
		log.Printf("• If sustained, will reach 2030 climate targets by 2045")
	} else {
		log.Printf("• ⚠️ Emissions increasing - requires policy intervention")
	}
}

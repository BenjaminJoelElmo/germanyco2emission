// process.go
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------
//
//	Data structure – keep the fields you need for the plot.
//	You can drop any you don’t use; the loader will simply leave them zero.
//
// ---------------------------------------------------------------------
type CO2Record struct {
	CountryName  string
	CountryCode  string
	Year         int
	CO2Emissions float64 // megatonnes CO₂e
	Temp         float64 // optional temperature column (°C)
}

// ---------------------------------------------------------------------
//
//	LoadAndClean – reads any CSV that contains a **Year** column
//	               and a **CO₂** (or “emission”) column. It works with
//	               – 5‑column “full” files (Country, Code, Year, CO₂, Temp)
//	               – 2‑column “minimal” files (Year, CO₂)
//	               – any superset of those (extra columns are ignored).
//
// ---------------------------------------------------------------------
func LoadAndClean(path string) ([]CO2Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	// ---------- 1. Read header ----------
	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}
	// Normalise header strings for easier matching
	for i := range header {
		header[i] = strings.TrimSpace(strings.ToLower(header[i]))
	}

	// ---------- 2. Detect column indexes ----------
	idxYear := -1
	idxCO2 := -1
	idxCountryName := -1
	idxCountryCode := -1
	idxTemp := -1

	for i, col := range header {
		switch {
		case strings.Contains(col, "year"):
			idxYear = i
		case strings.Contains(col, "co2") && strings.Contains(col, "emission"):
			idxCO2 = i
		case strings.Contains(col, "co2") && !strings.Contains(col, "emission"):
			// In a minimal file the column might just be called “co2”
			if idxCO2 == -1 {
				idxCO2 = i
			}
		case strings.Contains(col, "country") && strings.Contains(col, "name"):
			idxCountryName = i
		case strings.Contains(col, "country") && strings.Contains(col, "code"):
			idxCountryCode = i
		case strings.Contains(col, "temp"):
			idxTemp = i
		}
	}

	// If we didn’t find a year or CO₂ column, abort with a clear message.
	if idxYear == -1 || idxCO2 == -1 {
		return nil, fmt.Errorf("could not locate Year (idx=%d) or CO₂ column (idx=%d) in header: %v",
			idxYear, idxCO2, header)
	}

	// ---------- 3. Parse the rows ----------
	var records []CO2Record
	rowNum := 1 // we already consumed the header
	for {
		row, err := r.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("reading row %d: %w", rowNum, err)
		}
		rowNum++

		// Defensive check – a row may be shorter than the header
		if len(row) < idxYear+1 || len(row) < idxCO2+1 {
			log.Printf("⚠️ row %d has %d columns (expected ≥%d) – skipping",
				rowNum, len(row), max(idxYear, idxCO2)+1)
			continue
		}

		// ---------- Year ----------
		yrStr := strings.TrimSpace(row[idxYear])
		year, err := strconv.Atoi(yrStr)
		if err != nil {
			log.Printf("⚠️ row %d: invalid year '%s' – skipping", rowNum, yrStr)
			continue
		}

		// ---------- CO₂ ----------
		co2Str := strings.TrimSpace(row[idxCO2])
		co2, err := strconv.ParseFloat(co2Str, 64)
		if err != nil {
			log.Printf("⚠️ row %d: invalid CO₂ '%s' – setting to 0", rowNum, co2Str)
			co2 = 0
		}

		// ---------- Optional fields ----------
		countryName := ""
		if idxCountryName >= 0 && idxCountryName < len(row) {
			countryName = strings.TrimSpace(row[idxCountryName])
		}
		countryCode := ""
		if idxCountryCode >= 0 && idxCountryCode < len(row) {
			countryCode = strings.TrimSpace(row[idxCountryCode])
		}
		temp := 0.0
		if idxTemp >= 0 && idxTemp < len(row) {
			if tStr := strings.TrimSpace(row[idxTemp]); tStr != "" {
				if t, err := strconv.ParseFloat(tStr, 64); err == nil {
					temp = t
				}
			}
		}

		rec := CO2Record{
			CountryName:  countryName,
			CountryCode:  countryCode,
			Year:         year,
			CO2Emissions: co2,
			Temp:         temp,
		}
		records = append(records, rec)
	}
	return records, nil
}

// ---------------------------------------------------------------------
//
//	SaveCleaned – unchanged (keeps the same CSV format you already have)
//
// ---------------------------------------------------------------------
func SaveCleaned(records []CO2Record, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	// Header that matches the struct order (feel free to edit)
	if err := w.Write([]string{"CountryName", "CountryCode", "Year", "CO2Emissions", "Temp"}); err != nil {
		return err
	}
	for _, r := range records {
		row := []string{
			r.CountryName,
			r.CountryCode,
			strconv.Itoa(r.Year),
			fmt.Sprintf("%.3f", r.CO2Emissions),
			fmt.Sprintf("%.2f", r.Temp),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

// ---------------------------------------------------------------------
//
//	tiny helper – max of two ints
//
// ---------------------------------------------------------------------
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

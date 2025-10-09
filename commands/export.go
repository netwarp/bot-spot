package commands

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"main/database"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

func fileNamePrefix() string {
	// Create an "exports" folder if it doesn't exist
	if err := os.MkdirAll("exports", os.ModePerm); err != nil {
		panic(fmt.Errorf("failed to create exports folder: %w", err))
	}

	timestamp := time.Now().Format("2006-01-02 15-04-05")

	filePrefix := fmt.Sprintf("exports/%s", timestamp)
	return filePrefix
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func ToCSV(displayLogs bool) {

	fileName := fileNamePrefix() + ".csv"
	if displayLogs {
		color.Yellow("Export data to CSV file: " + fileName)
	}

	file, err := os.Create(RootDir() + "/" + filepath.Clean(fileName))
	if err != nil {
		panic(fmt.Errorf("failed to create file: %w", err))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close file: %w", err))
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Id",
		"Exchange",
		"Status",
		"Quantity",
		"BuyPrice",
		"SellPrice",
		"Gain USD",
		"BuyId",
		"SellId",
	}
	if err := writer.Write(header); err != nil {
		panic(fmt.Errorf("failed to write header: %w", err))
	}

	// Write each row
	cycles, err := database.CycleList()
	if err != nil {
		panic(fmt.Errorf("failed to get cycles: %w", err))
	}

	for _, cycle := range cycles {
		row := []string{
			fmt.Sprintf("%v", cycle.Id),
			fmt.Sprintf("%v", cycle.Exchange),
			fmt.Sprintf("%v", cycle.Status),
			fmt.Sprintf("%v", cycle.Quantity),
			fmt.Sprintf("%v", cycle.Buy.Price),
			fmt.Sprintf("%v", cycle.Sell.Price),
			fmt.Sprintf("%v", CalcAbsoluteGainByCycle(&cycle)),
			fmt.Sprintf("%v", cycle.Buy.ID),
			fmt.Sprintf("%v", cycle.Sell.ID),
		}

		if err := writer.Write(row); err != nil {
			panic(fmt.Errorf("failed to write row: %w", err))
		}
	}
	if displayLogs {
		color.Green("Successfully Export data to CSV file: " + fileName)
	}
}

func ToJSON(displayLogs bool) {
	fileName := fileNamePrefix() + ".json"
	fileName = RootDir() + "/" + filepath.Clean(fileName)

	cycles, err := database.CycleList()
	if err != nil {
		panic(fmt.Errorf("failed to get database: %w", err))
	}

	jsonData, err := json.MarshalIndent(cycles, "", "  ")
	if err != nil {
		panic(fmt.Errorf("failed to marshal json: %w", err))
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write JSON file '%s': %w", fileName, err))
	}

	if displayLogs {
		color.Green("Successfully Export data to JSON file: " + fileName)
	}
}

func Export(displayLogs bool) {
	ToCSV(displayLogs)
	ToJSON(displayLogs)
}

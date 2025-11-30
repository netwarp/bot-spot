package commands

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"main/database"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

func ensureExportsDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("error getting home directory: %v", err))
	}

	exportsDir := filepath.Join(homeDir, "exports")
	err = os.MkdirAll(exportsDir, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("error creating exports directory: %v", err))
	}

	return exportsDir
}

func filePrefix() string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	completePath := filepath.Join(ensureExportsDir(), timestamp)
	return completePath
}

func toJSON() {
	file := filePrefix() + ".json"

	fmt.Println("Exporting to:", file)

	cycles, err := database.CycleList()
	if err != nil {
		panic(fmt.Errorf("error getting cycles: %v", err))
	}

	jsonData, err := json.MarshalIndent(cycles, "", "  ")
	if err != nil {
		panic(fmt.Errorf("error marshalling cycles to JSON: %v", err))
	}

	err = os.WriteFile(file, jsonData, 0644)
	if err != nil {
		panic(fmt.Errorf("error writing file: %v", err))
	}

	color.Green(file)
}

func toCSV() {
	fileName := filePrefix() + ".csv"

	file, err := os.Create(fileName)
	if err != nil {
		panic(fmt.Errorf("error creating file: %v", err))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(fmt.Errorf("error closing file: %v", err))
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
		"BuyId",
		"SellId",
		"Free balance",
		"Dedicated balance",
		"Buy offset",
		"Sell offset",
		"Percent",
		"BTC price",
		"Absolute gain",
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
			fmt.Sprintf("%v", cycle.Buy.ID),
			fmt.Sprintf("%v", cycle.Sell.ID),
			fmt.Sprintf("%v", cycle.MetaData.FreeBalanceUSD),
			fmt.Sprintf("%v", cycle.MetaData.USDDedicated),
			fmt.Sprintf("%v", cycle.Buy.Offset),
			fmt.Sprintf("%v", cycle.Sell.Offset),
			fmt.Sprintf("%v", cycle.MetaData.Percent),
			fmt.Sprintf("%v", cycle.MetaData.BTCPrice),
			fmt.Sprintf("%v", cycle.CalcProfit()),
		}

		if err := writer.Write(row); err != nil {
			panic(fmt.Errorf("failed to write row: %w", err))
		}
	}

	color.Green(fileName)
}

func Export() {
	toCSV()
	toJSON()
}

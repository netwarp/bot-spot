package tools

import (
	"encoding/json"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"log"
	"main/database"
	"os"
	"path/filepath"
)

type CycleOldModel struct {
	ID        string  `json:"_id"`
	BuyID     string  `json:"buyId"`
	BuyPrice  float64 `json:"buyPrice"`
	Exchange  string  `json:"exchange"`
	IDInt     int     `json:"idInt"`
	Quantity  float64 `json:"quantity"`
	SellID    string  `json:"sellId"`
	SellPrice float64 `json:"sellPrice"`
	Status    string  `json:"status"`
}

// FromCloverToSqlite is reserved for future migration tools.
func FromCloverToSqlite() {
	err := godotenv.Load("../bot.conf")
	if err != nil {
		log.Fatal("Error loading ../bot.conf")
	}

	exportFile := os.Getenv("EXPORT_FILE")
	if exportFile == "" {
		log.Fatal("Missing environment variable: EXPORT_FILE (e.g., EXPORT_FILE=2025-10-12 16-39-54.json)")
	}

	filePath := filepath.Join("../exports", exportFile)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(fileContent)

	var cycleOldModel []CycleOldModel
	err = json.Unmarshal([]byte(jsonString), &cycleOldModel)
	if err != nil {
		log.Fatalf("Error unmarshalling orders JSON: %v", err)
	}

	log.Printf("Successfully loaded %d orders from %s.", len(cycleOldModel), exportFile)

	for _, oldCycle := range cycleOldModel {
		id := oldCycle.IDInt
		status := oldCycle.Status
		quantity := oldCycle.Quantity
		buyPrice := oldCycle.BuyPrice
		sellPrice := oldCycle.SellPrice
		exchange := oldCycle.Exchange

		var cycle database.Cycle
		cycle.Id = id
		cycle.Status = database.Status(status)
		cycle.Quantity = quantity
		cycle.Buy.Price = buyPrice
		cycle.Sell.Price = sellPrice
		cycle.Exchange = exchange
		cycle.Buy.ID = oldCycle.BuyID
		cycle.Sell.ID = oldCycle.SellID

		_, err := database.CycleNew(&cycle)
		if err != nil {
			log.Fatal(err)
		}
	}

	color.Green("Successfully migrated %d orders from %s.", len(cycleOldModel), exportFile)
}

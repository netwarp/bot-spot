package tools

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

// FromCloverToSqlite is reserved for future migration tools.
func FromCloverToSqlite() {
	err := godotenv.Load("../bot.conf")
	if err != nil {
		log.Fatal("Error loading ../bot.conf")
	}

	exportFile := os.Getenv("EXPORT_FILE")
	if exportFile == "" {
		log.Fatal("Missing environment variable: EXPORT_FILE=2025-10-12 16-39-54.json")
	}

	filePath := filepath.Join("../exports", exportFile)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(fileContent)

	var prettyJSON interface{}
	err = json.Unmarshal([]byte(jsonString), &prettyJSON)
	if err != nil {
		log.Fatalf("Error unmarshalling orders JSON: %v", err)
	}
}

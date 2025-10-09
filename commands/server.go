package commands

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"html/template"
	"log"
	"main/database"
	"net/http"
	"os"
	"strconv"
)

var perPage int = 200

func getPage(r *http.Request) int {
	query := r.URL.Query()
	pageStr := query.Get("page")

	if pageStr == "" {
		return 1
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}

	return page
}

func getAddressServer() string {
	if os.Getenv("SERVER_ADDRESS") != "" {
		return os.Getenv("SERVER_ADDRESS")
	}
	return "localhost:8080"
}

func Server() {
	MainMiddleware()

	var address = getAddressServer()

	fmt.Println("Open browser then go to " + address)
	color.Magenta("\nCtrl + C to close the server")
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := getPage(r)
		fmt.Println(page)

		//docs := database.ListPerPage(page, perPage)
		cycles, err := database.CycleList()
		if err != nil {
			http.Error(w, "Error getting cycles", http.StatusInternalServerError)
			return
		}

		cyclesCount := 0
		cyclesCompleted := 0
		totalBuy := 0.0
		totalSell := 0.0

		for _, _ = range cycles {
			//quantity := cycle.Quantity
			//buyPrice := cycle.Buy.Price
			//sellPrice := cycle.Sell.Price
			//status := cycle.Status

			//var quantityFloat, buyPriceFloat, sellPriceFloat float64
			//var quantityStr string
			//
			//var percentageChange string
			//if buyPriceFloat > 0 {
			//	change := ((quantityFloat * sellPriceFloat) - (quantityFloat * buyPriceFloat)) / (quantityFloat * buyPriceFloat) * 100
			//	percentageChange = fmt.Sprintf("%.2f%%", change)
			//} else {
			//	percentageChange = "N/A"
			//}

			// gain usd
			//gainUSD := (sellPriceFloat - buyPriceFloat) * quantityFloat

			// Update stats
			cyclesCount++

			//if status == "completed" {
			//	cyclesCompleted++
			//
			//	b := (buyPrice).(float64) * (quantity).(float64)
			//	totalBuy += b
			//
			//	s := (sellPrice).(float64) * (quantity).(float64)
			//	totalSell += s
			//}
		}

		gainAbs := totalSell - totalBuy
		percent := (totalSell - totalBuy) / totalBuy * 100
		// Pagination

		tmpl, err := template.ParseFiles("commands/misc/template.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, map[string]interface{}{
			"Cycles":          cycles,
			"cyclesCount":     cyclesCount,
			"cyclesCompleted": cyclesCompleted,
			"totalBuy":        totalBuy,
			"totalSell":       totalSell,
			"gainAbs":         gainAbs,
			"page":            page,
			"percent":         percent,
		})
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/api/get-order", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			OrderID  string `json:"orderId"`
			Exchange string `json:"exchange"`
		}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			http.Error(w, `{"error": "invalid json"}`, http.StatusBadRequest)
			return
		}

		client := GetClientByExchange(data.Exchange)
		order, err := client.GetOrderById(data.OrderID)
		if err != nil {
			http.Error(w, `{"error": "order not found"}`, http.StatusNotFound)
			return
		}

		_, err = w.Write(order)
		if err != nil {
			return
		}

	})

	err := http.ListenAndServe(address, mux)
	if err != nil {
		log.Fatal(err)
	}
}

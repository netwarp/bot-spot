package commands

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"html/template"
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

func Server() error {
	MainMiddleware()

	var address = getAddressServer()

	fmt.Println("Open browser then go to " + address)
	color.Magenta("\nCtrl + C to close the server")

	mux := http.NewServeMux()

	mux.HandleFunc("/", displayStats)

	mux.HandleFunc("/api/get-order", getOrder)

	err := http.ListenAndServe(address, mux)
	if err != nil {
		return fmt.Errorf("error starting server: %v", err)
	}

	return nil
}

func displayStats(w http.ResponseWriter, r *http.Request) {
	page := getPage(r)
	fmt.Println("Page: ", page)

	cycles, err := database.CycleList()

	if err != nil {
		http.Error(w, "Error getting cycles", http.StatusInternalServerError)
		return
	}

	cyclesCount := 0
	cyclesCompleted := 0
	totalBuy := 0.0
	totalSell := 0.0
	totalProfit := 0.0

	for _, cycle := range cycles {
		//fmt.Printf("%+v\n", cycle)
		cyclesCount++
		if cycle.Status == database.Completed {
			cyclesCompleted++

			totalBuy += cycle.Buy.Price * cycle.Quantity
			totalSell += cycle.Sell.Price * cycle.Quantity

			totalProfit += cycle.CalcProfit()
		}

	}

	tmpl, err := template.ParseFiles("commands/misc/template.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, map[string]interface{}{
		"cycles":          cycles,
		"cyclesCount":     cyclesCount,
		"cyclesCompleted": cyclesCompleted,
		"totalBuy":        totalBuy,
		"totalSell":       totalSell,
		"totalProfit":     totalProfit,
		"page":            page,
	})

	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func getOrder(w http.ResponseWriter, r *http.Request) {
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
}

package commands

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"log"
	"main/database"
	"main/tools"
	"math"
	"os"
	"strconv"
)

func CalcAmountUSD(freeBalance float64, percentStr string) float64 {
	percent, err := strconv.ParseFloat(percentStr, 64)
	if err != nil {
		log.Fatal(err)
	}
	return percent * freeBalance / 100
}

func CalcAmountBTC(availableUSD, priceBTC float64) float64 {
	return availableUSD / priceBTC
}

func FormatSmallFloat(quantity float64) string {
	return fmt.Sprintf("%.6f", quantity)
}

func New() error {
	MainMiddleware()

	percent := os.Getenv("PERCENT")

	buyOffset, _ := strconv.ParseFloat(os.Getenv("BUY_OFFSET"), 64)
	buyOffset = math.Abs(buyOffset)

	sellOffset, _ := strconv.ParseFloat(os.Getenv("SELL_OFFSET"), 64)
	//sellOffset = math.Abs(sellOffset)

	client := GetClientByExchange()

	client.CheckConnection()

	freeBalance := client.GetBalanceUSD()
	color.Cyan("Free USD Balance: %.2f", freeBalance)
	if freeBalance < 10 {
		tools.Telegram("At least 10$ needed")
		color.Red("At least 10$ needed")
		os.Exit(0)
	}

	btcPrice := client.GetLastPriceBTC()

	fmt.Printf("%s %s\n",
		color.CyanString("BTC Price"),
		color.YellowString("%.2f", btcPrice),
	)

	newCycleUSDC := CalcAmountUSD(freeBalance, percent)

	fmt.Printf("%s %s\n",
		color.CyanString("USD for this new cycle:"),
		color.YellowString("%.2f", newCycleUSDC),
	)

	newCycleBTC := CalcAmountBTC(newCycleUSDC, btcPrice)
	newCycleBTCFormated := FormatSmallFloat(newCycleBTC)
	fmt.Printf("%s %s\n",
		color.CyanString("BTC for this new cycle:"),
		color.YellowString(newCycleBTCFormated),
	)

	buyPrice := btcPrice - buyOffset
	fmt.Printf("%s %s\n",
		color.CyanString("Buy Price"),
		color.YellowString("%.2f", buyPrice),
	)

	sellPrice := btcPrice + sellOffset
	fmt.Printf("%s %s\n",
		color.CyanString("Sell Price"),
		color.YellowString("%.2f", sellPrice),
	)

	// Prepare Order
	buyPriceStr := fmt.Sprintf("%.2f", buyPrice)

	body, err := client.CreateOrder("BUY", buyPriceStr, newCycleBTCFormated)
	if err != nil {
		color.Red("Order failed:", err)
		tools.Telegram("Order failed: " + err.Error())
		os.Exit(0)
	}

	orderId, _, _, err := jsonparser.Get(body, "orderId")
	if err != nil {
		tools.Telegram("Order failed: " + err.Error())
		log.Fatal(err)
	}

	// Insert in database
	cycle := database.Cycle{
		Exchange:  "mexc", // Todo dynamic
		Status:    "buy",
		Quantity:  newCycleBTC,
		BuyPrice:  buyPrice,
		BuyId:     string(orderId),
		SellPrice: sellPrice,
		SellId:    "",
	}
	docId := database.NewCycle(&cycle)

	message := "New Cycle successfully inserted in database"
	color.Green(message)
	Log(message)

	// Save data without logs
	Export(false)

	if os.Getenv("TELEGRAM") == "1" {
		doc := database.GetById(docId)
		idInt := doc.Get("idInt")

		var message = ""
		message += fmt.Sprintf("ℹ️ New Cycle: %d \n", idInt)
		message += fmt.Sprintf("✨ Quantity: %.6f \n", newCycleBTC)
		message += fmt.Sprintf("📉 Buy Price: %.2f \n", buyPrice)
		message += fmt.Sprintf("📈 Sell Price: %.2f \n", sellPrice)
		tools.Telegram(message)
	}
}

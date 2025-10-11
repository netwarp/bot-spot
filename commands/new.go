package commands

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"log"
	"main/database"
	"main/tools"
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

func New() error {
	MainMiddleware()

	newCycle, err := PrepareNewCycle()
	if err != nil {
		return fmt.Errorf("error preparing new cycle: %v", err)
	}

	client := GetClientByExchange(newCycle.Exchange)

	//// Prepare Order
	buyPriceStr := fmt.Sprintf("%.2f", newCycle.Buy.Price)
	buyBTCQuantityStr := fmt.Sprintf("%.6f", newCycle.Quantity)

	body, err := client.CreateOrder("BUY", buyPriceStr, buyBTCQuantityStr)
	if err != nil {
		color.Red("Order failed:", err)
		tools.Telegram("Order failed: " + err.Error())
		os.Exit(0)
	}
	buyOrderId, _, _, err := jsonparser.Get(body, "orderId")
	if err != nil {
		tools.Telegram("Order failed: " + err.Error())
		log.Fatal("Order failed: " + err.Error())
	}

	newCycle.Buy.ID = string(buyOrderId)
	newCycle.Status = database.Buy

	// Insert in database
	_, err = database.CycleNew(newCycle)
	if err != nil {
		return fmt.Errorf("error inserting new cycle in database: %v", err)
	}
	message := "New Cycle successfully inserted in database"
	color.Green(message)
	Log(message)

	notifTelegram(newCycle)

	return nil
}

// PrepareNewCycle Prepare new cycle before place order and insert in db
func PrepareNewCycle() (*database.Cycle, error) {
	newCycle := database.Cycle{}

	// Exchange
	exchange := getExchange()
	newCycle.Exchange = exchange

	// Percent
	percent := getPercent()
	newCycle.MetaData.Percent = percent

	// BuyOffset
	buyOffset := getOffset("BUY_OFFSET")
	newCycle.Buy.Offset = buyOffset

	// SellOffset
	sellOffset := getOffset("SELL_OFFSET")
	newCycle.Sell.Offset = sellOffset

	client := GetClientByExchange(exchange)
	client.CheckConnection()

	// BTCPrice
	btcPrice, err := client.GetLastPriceBTC()
	if err != nil {
		return nil, err
	}
	newCycle.MetaData.BTCPrice = btcPrice

	// BuyPrice
	buyPrice := btcPrice + float64(newCycle.Buy.Offset)
	newCycle.Buy.Price = buyPrice

	// Sell Price
	sellPrice := btcPrice + float64(newCycle.Sell.Offset)
	newCycle.Sell.Price = sellPrice

	// FreeBalanceUSD
	freeBalance, err := client.GetBalanceUSD()
	if err != nil {
		return nil, fmt.Errorf("error getting free balance: %v", err)
	}
	if freeBalance < 10 {
		color.Red("At least 10$ needed")
		os.Exit(0)
	}
	newCycle.MetaData.FreeBalanceUSD = freeBalance

	// USDDedicated
	usdDedicated := CalcAmountUSD(freeBalance, strconv.Itoa(newCycle.MetaData.Percent))
	newCycle.MetaData.USDDedicated = usdDedicated

	// BTCQuantity
	btcQuantity := CalcAmountBTC(newCycle.MetaData.USDDedicated, newCycle.Buy.Price)
	newCycle.Quantity = btcQuantity

	// Display Data
	const fieldWidth = 27
	var formatString = "%-" + strconv.Itoa(fieldWidth) + "s %s\n"

	fmt.Printf(formatString,
		color.CyanString("Exchange"),
		color.YellowString(newCycle.Exchange),
	)

	fmt.Printf(formatString,
		color.CyanString("Percent"),
		color.YellowString(strconv.Itoa(newCycle.MetaData.Percent)),
	)

	fmt.Printf(formatString,
		color.CyanString("Buy Offset"),
		color.YellowString(strconv.Itoa(newCycle.Buy.Offset)),
	)

	fmt.Printf(formatString,
		color.CyanString("Sell Offset"),
		color.YellowString(strconv.Itoa(newCycle.Sell.Offset)),
	)

	fmt.Printf(formatString,
		color.CyanString("BTC Price"),
		color.YellowString("%.2f", newCycle.MetaData.BTCPrice),
	)

	fmt.Printf(formatString,
		color.CyanString("Buy Price"),
		color.YellowString("%.2f", newCycle.Buy.Price),
	)

	fmt.Printf(formatString,
		color.CyanString("Sell Price"),
		color.YellowString("%.2f", newCycle.Sell.Price),
	)

	fmt.Printf(formatString,
		color.CyanString("Free Balance"),
		color.YellowString("%.2f", newCycle.MetaData.FreeBalanceUSD),
	)

	fmt.Printf(formatString,
		color.CyanString("Dedicated Balance"),
		color.YellowString("%.2f", newCycle.MetaData.USDDedicated),
	)

	fmt.Printf(formatString,
		color.CyanString("BTC Quantity"),
		color.YellowString("%.6f", newCycle.Quantity),
	)

	return &newCycle, nil
}

func getExchange() string {
	exchange := os.Getenv("EXCHANGE")
	if exchange == "" {
		color.Red("EXCHANGE env variable is required")
		os.Exit(0)
	}
	return exchange
}

func getPercent() int {
	percentStr := os.Getenv("PERCENT")
	if percentStr == "" {
		color.Red("PERCENT env variable is required")
		os.Exit(0)
	}
	percent, err := strconv.Atoi(percentStr)
	if err != nil {
		color.Red("PERCENT env variable must be a number")
		os.Exit(0)
	}
	return percent
}

func getOffset(key string) int {
	offset := os.Getenv(key)
	if offset == "" {
		color.Red(key + " env variable is required")
		os.Exit(0)
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		color.Red(key + " env variable must be a number")
		os.Exit(0)
	}
	return offsetInt
}

func notifTelegram(cycle *database.Cycle) {
	if os.Getenv("TELEGRAM") == "1" {
		var message = ""
		message += fmt.Sprintf("ℹ️ New Cycle: %d \n", cycle.Id)
		message += fmt.Sprintf("✨ Quantity: %.6f \n", cycle.Quantity)
		message += fmt.Sprintf("📉 Buy Price: %.2f \n", cycle.Buy.Price)
		message += fmt.Sprintf("📈 Sell Price: %.2f \n", cycle.Sell.Price)
		tools.Telegram(message)
	}
}

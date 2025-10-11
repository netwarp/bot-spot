package commands

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"log"
	"main/database"
	"main/tools"
	"strconv"
)

var client ExchangeClient = nil
var lastPrice float64 = 0.0

func Update() error {
	MainMiddleware()

	client = GetClientByExchange()
	client.CheckConnection()

	lastPrice, _ = client.GetLastPriceBTC()
	fmt.Println("Last price:", lastPrice)

	cycles, err := database.CycleList()
	if err != nil {
		return fmt.Errorf("error getting cycles: %v", err)
	}

	for _, cycle := range cycles {
		if cycle.Status == "buy" {
			err := handleBuy(&cycle)
			if err != nil {
				return fmt.Errorf("error handling buy: %v", err)
			}
		} else if cycle.Status == "sell" {
			err := handleSell(&cycle)
			if err != nil {
				return fmt.Errorf("error handling sell: %v", err)
			}
		}
	}

	return nil
}

func handleBuy(cycle *database.Cycle) error {
	buyOrderId := cycle.Buy.ID

	order, err := client.GetOrderById(buyOrderId)
	if err != nil {
		return fmt.Errorf("error getting order: %v", err)
	}

	isFilled, err := client.IsFilled(string(order))
	if err != nil {
		return fmt.Errorf("error checking order: %v", err)
	}

	if !isFilled {
		fmt.Printf("%s %s %s\n",
			color.YellowString("%d", cycle.Id),
			color.CyanString("Order Buy still active -"),
			color.WhiteString("%s", buyOrderId),
		)

		return nil
	}

	fmt.Printf("%s %s\n",
		color.YellowString("%d", cycle.Id),
		color.GreenString("Order Buy filled"),
	)

	sellPrice := cycle.Sell.Price

	if lastPrice > cycle.Sell.Price {
		upOffset := 200.0
		newSellPrice := cycle.Sell.Price + upOffset
		sellPrice = newSellPrice
		fmt.Println("New sell price: ", newSellPrice)

		_, err := database.CycleUpdate(cycle.Id, "sellPrice", newSellPrice)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("New sell price updated: ")
	}

	quantity := cycle.Quantity
	quantityStr := strconv.FormatFloat(quantity, 'f', 6, 64)
	sellPriceStr := strconv.FormatFloat(sellPrice, 'f', 6, 64)

	bytes, err := client.CreateOrder("SELL", sellPriceStr, quantityStr)
	if err != nil {
		return fmt.Errorf("error creating sell order: %v", err)
	}

	orderId, _, _, err := jsonparser.Get(bytes, "orderId")
	if err != nil {
		log.Printf("Cycle id %d", cycle.Id)
		log.Printf("Failed to parse orderId: %v", err)
		log.Fatal(err)
	}

	fmt.Printf("%s %s %s\n",
		color.YellowString("%d", cycle.Id),
		color.CyanString("New sell Order -"),
		color.WhiteString("%s", string(bytes)),
	)

	_, err = database.CycleUpdate(cycle.Id, "status", database.Sell)
	if err != nil {
		return fmt.Errorf("error updating cycle status: %v", err)
	}
	_, err = database.CycleUpdate(cycle.Id, "sellId", string(orderId))
	if err != nil {
		return fmt.Errorf("error updating cycle sell id: %v", err)
	}

	return nil
}

func handleSell(cycle *database.Cycle) error {
	sellOrderId := cycle.Sell.ID
	order, err := client.GetOrderById(sellOrderId)
	if err != nil {
		return fmt.Errorf("error getting order: %v", err)
	}

	isFilled, err := client.IsFilled(string(order))
	if err != nil {
		return fmt.Errorf("error checking order: %v", err)
	}

	if !isFilled {
		fmt.Printf("%s %s %s\n",
			color.YellowString("%d", cycle.Id),
			color.CyanString("Order Sell still active -"),
			color.WhiteString("%s", sellOrderId),
		)

		return nil
	}

	fmt.Printf("%s %s\n",
		color.YellowString("%d", cycle.Id),
		color.GreenString("Order Sell filled"),
	)

	_, err = database.CycleUpdate(cycle.Id, "status", database.Completed)
	if err != nil {
		return fmt.Errorf("error updating cycle status: %v", err)
	}

	color.Green("Cycle successfully completed")

	notifTelegram2(cycle)

	return nil
}

func notifTelegram2(cycle *database.Cycle) {
	var message = ""
	message += fmt.Sprintf("✅ Cycle %d completed \n", cycle.Id)
	message += fmt.Sprintf("📉 Buy Price: %.2f \n", cycle.Buy.Price)
	message += fmt.Sprintf("📈 Sell Price: %.2f \n", cycle.Sell.Price)
	message += fmt.Sprintf("💰 Gain: $ %.2f \n", cycle.CalcProfit())
	tools.Telegram(message)
}

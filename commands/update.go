package commands

func Update() error {
	//MainMiddleware()

	//client := GetClientByExchange()
	//client.CheckConnection()
	//
	//lastPrice := client.GetLastPriceBTC()
	////fmt.Println("Last price:", lastPrice)
	//
	//cycles, err := database.CycleList()
	//if err != nil {
	//	return fmt.Errorf("error getting cycles: %v", err)
	//}
	//
	//for _, cycle := range cycles {
	//
	//	id := cycle.Id
	//	status := cycle.Status
	//	quantity := cycle.Quantity
	//	buyPrice := cycle.Buy.Price
	//	buyId := cycle.Buy.ID
	//	sellPrice := cycle.Sell.Price
	//	sellId := cycle.Sell.ID
	//
	//	if status == "buy" {
	//		order, _ := client.GetOrderById(buyId)
	//		isFilled, err := client.IsFilled(string(order))
	//		if err != nil {
	//			log.Println(err.Error())
	//			color.Red("Error found on cycle %d (don't worry)\n", id)
	//			color.Yellow("Try to remove it")
	//
	//			fmt.Printf("go run . -c %d\n", id)
	//			fmt.Printf("go run . -cl %d %d\n", id, id)
	//			os.Exit(0)
	//		}
	//
	//		if !isFilled {
	//			fmt.Printf("%s %s %s\n",
	//				color.YellowString("%d", id),
	//				color.CyanString("Order Buy  still active -"),
	//				color.WhiteString("%s", buyId),
	//			)
	//		} else {
	//			fmt.Printf("%s %s\n",
	//				color.YellowString("%d", id),
	//				color.GreenString("Order Buy filled"),
	//			)
	//
	//			// Check sell price > last price
	//			fmt.Println("SELL PRICE", sellPrice)
	//			fmt.Println("id ", id)
	//			sellPrice := sellPrice
	//
	//			if lastPrice > sellPrice {
	//				const offset = 200.0
	//				newSellPrice := sellPrice + offset
	//				fmt.Println("New sell price: ", newSellPrice)
	//
	//				_, err = database.CycleUpdate(id, "sellPrice", newSellPrice)
	//				if err != nil {
	//					log.Fatal(err)
	//				}
	//				fmt.Println("New sell price updated: ")
	//			}
	//
	//			// Place sell order
	//			quantity := (quantity).(float64)
	//			quantityStr := strconv.FormatFloat(quantity, 'f', 6, 64)
	//
	//			doc := database.GetById(idString)
	//			sellPrice = doc.Get("sellPrice").(float64)
	//			sellPriceStr := strconv.FormatFloat(sellPrice, 'f', 6, 64)
	//
	//			bytes, err := client.CreateOrder("SELL", sellPriceStr, quantityStr)
	//			if err != nil {
	//				log.Fatalf("Error creating order: %v", err)
	//				return
	//			}
	//			orderId, _, _, err := jsonparser.Get(bytes, "orderId")
	//
	//			if err != nil {
	//				log.Printf("Cycle id %d", idInt)
	//				log.Printf("Failed to parse orderId: %v", err)
	//				log.Fatal(err)
	//			}
	//
	//			fmt.Printf("%s %s %s\n",
	//				color.YellowString("%d", idInt),
	//				color.CyanString("New sell Order -"),
	//				color.WhiteString("%s", string(bytes)),
	//			)
	//
	//			database.FindCycleByIdAndUpdate(idString, "status", "sell")
	//			database.FindCycleByIdAndUpdate(idString, "sellId", string(orderId))
	//		}
	//	} else if status == "sell" {
	//		order, _ := client.GetOrderById((sellId).(string))
	//		isFilled, err := client.IsFilled(string(order))
	//		if err != nil {
	//			color.Red("Error found on cycle %d (don't worry)\n", idInt)
	//			log.Println(err)
	//			color.Yellow("Try to remove it")
	//
	//			fmt.Printf("go run . -c %d\n", idInt)
	//			fmt.Printf("go run . -cl %d %d\n", idInt, idInt)
	//			os.Exit(0)
	//		}
	//
	//		if !isFilled {
	//			fmt.Printf("%s %s %s\n",
	//				color.YellowString("%d", idInt),
	//				color.CyanString("Order Sell still active -"),
	//				color.WhiteString("%s", sellId),
	//			)
	//		} else {
	//			database.FindCycleByIdAndUpdate(idString, "status", "completed")
	//
	//			// Calc Percent
	//			totalBuyUSD := buyPrice.(float64) * quantity.(float64)
	//			totalSellUSD := sellPrice.(float64) * quantity.(float64)
	//			percent := (totalSellUSD - totalBuyUSD) / totalBuyUSD * 100
	//
	//			fmt.Printf("%s %s (Gain: %.2f%%)\n",
	//				color.YellowString("%d", idInt),
	//				color.GreenString("Cycle successfully completed"),
	//				percent,
	//			)
	//
	//			if os.Getenv("TELEGRAM") == "1" {
	//				var message = ""
	//				message += fmt.Sprintf("✅ Cycle %d completed \n", idInt)
	//				message += fmt.Sprintf("📉 Buy Price: %.2f \n", buyPrice)
	//				message += fmt.Sprintf("📈 Sell Price: %.2f \n", sellPrice)
	//				message += fmt.Sprintf("💰 Gain: $ %.2f \n", totalSellUSD-totalBuyUSD)
	//				tools.Telegram(message)
	//			}
	//		}
	//	}
	//
	//	// Add a sleeper to ensure lock/unlock works fine
	//	time.Sleep(50 * time.Millisecond)
	//}
	//Log("Update complete")
	//Export(false)

	return nil
}

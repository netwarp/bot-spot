package commands

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"main/database"
	"os"
	"strconv"
)

func Cancel() error {
	if len(os.Args) < 3 {
		color.Red("Id required")
		color.Cyan("go run . -c 34")
		os.Exit(1)
	}

	lastArg := os.Args[2]

	id, err := strconv.Atoi(lastArg)
	if err != nil {
		return err
	}

	color.Yellow("Cancelling %d", id)

	cycle, err := database.CycleGetById(id)
	if err != nil {
		return err
	}

	status := cycle.Status
	if status == "completed" {
		errMsg := "can't cancel completed cycle, only 'buy' or 'sell' is supported"
		color.Red(errMsg)
		//os.Exit(0)
		return fmt.Errorf(errMsg)
	}

	client := GetClientByExchange()

	buyId := cycle.Buy.ID
	sellId := cycle.Sell.ID

	res, err := client.CancelOrder(buyId)
	if err != nil {
		log.Println(string(res))
		return err
	}
	fmt.Println(string(res))

	res, err = client.CancelOrder(sellId)
	if err != nil {
		log.Println(string(res))
		return err
	}
	fmt.Println(string(res))

	err = database.CycleDeleteById(id)
	if err != nil {
		return err
	}

	color.Green("Cycle %d successfully canceled", id)

	return nil
}

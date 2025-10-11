package mexc

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

var client *Client

func TestMain(m *testing.M) {
	// TODO make config folder
	_ = godotenv.Load("../../bot.conf")

	client = NewClient()
	client.SetBaseURL("https://api.mexc.co")

	os.Exit(m.Run())
}

func TestCheckConnection(t *testing.T) {
	client.CheckConnection()
}

func TestClient_GetBalanceUSD(t *testing.T) {
	balance, _ := client.GetBalanceUSD()
	fmt.Println("balance:", balance)
}

func TestClient_GetOrderById(t *testing.T) {
	orderId := os.Getenv("ORDER_ID")
	order, err := client.GetOrderById(orderId)

	if err != nil {
		t.Error(err)
	}

	orderJSON, err := jsonparser.ParseString(order)

	fmt.Println(orderJSON)

	isFilled, _ := client.IsFilled(orderJSON)
	fmt.Println("is filled:", isFilled)
}

func TestCreateOrder(t *testing.T) {
	side := "SELL"
	price := "86600"
	quantity := "0.000025"

	order, err := client.CreateOrder(side, price, quantity)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(order))
}

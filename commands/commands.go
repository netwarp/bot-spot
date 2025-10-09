package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"io"
	"log"
	"main/database"
	"main/exchanges/mexc"
	"net/http"
	"os"
	"strings"
	"time"
)

type ExchangeClient interface {
	CheckConnection()
	GetBalanceUSD() (float64, error)
	GetLastPriceBTC() (float64, error)
	SetBaseURL(url string)
	CreateOrder(side, price, quantity string) ([]byte, error)
	GetOrderById(id string) ([]byte, error)
	IsFilled(id string) (bool, error)
	CancelOrder(orderID string) ([]byte, error)
	GetOpenOrders() ([]byte, error)
}

const ConfigFilename = "bot.conf"

func CreateConfigFileIfNotExists() {
	if _, err := os.Stat(ConfigFilename); errors.Is(err, os.ErrNotExist) {
		pathConfTemplate := fmt.Sprintf("commands/misc/%s.example", ConfigFilename)

		content, err := os.ReadFile(pathConfTemplate)
		if err != nil {
			content = []byte("CUSTOMER_ID=\n\nEXCHANGE=\n\nMEXC_PUBLIC=\nMEXC_PRIVATE=\n\nBUY_OFFSET=-1000\nSELL_OFFSET=1000\n\nPERCENT=6\n\nAUTO_INTERVAL_NEW=60\nAUTO_INTERVAL_UPDATE=10")
			err := os.WriteFile(ConfigFilename, content, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = os.WriteFile(ConfigFilename, content, 0644)
		if err != nil {
			log.Fatal(err)
		}
		color.Green("Config file created: " + ConfigFilename)
	}

}

func LoadDotEnv() {
	err := godotenv.Load(ConfigFilename)
	if err != nil {
		log.Fatal("Error loading bot.conf")
	}
}

func MainMiddleware() {
	color.Blue("Checking subscription before...")
	var customerId string = os.Getenv("CUSTOMER_ID")

	if customerId == "" {
		color.Red("You need to set CUSTOMER_ID in bot.conf")
		os.Exit(0)
	}

	url := "https://validator.cryptomancien.com"
	body, err := json.Marshal(
		map[string]string{
			"CUSTOMER_ID": os.Getenv("CUSTOMER_ID"),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	statusCode := resp.StatusCode
	if statusCode != 200 {
		color.Red("CUSTOMER_ID not in bot.conf or subscription expired")
		color.Red("Go to https://cryptomancien.com -> Space -> Trading bots and get your CUSTOMER_ID")
		color.Red("Then fill it in your config file bot.conf")
		os.Exit(0)
	}

	color.Green("Subscription OK")
	fmt.Println("")
}

func GetClientByExchange(exchangeArg ...string) ExchangeClient {

	var ex string
	if len(exchangeArg) > 0 {
		ex = exchangeArg[0]
	} else {
		ex = os.Getenv("EXCHANGE")
	}
	ex = strings.ToUpper(ex)

	var client ExchangeClient

	switch ex {
	case "MEXC":
		client = mexc.NewClient()
		client.SetBaseURL("https://api.mexc.co")
	default:
		fmt.Println("Unsupported exchange:", ex)
		os.Exit(0)
	}

	return client
}

func Log(message string) {
	logDir := "logs"

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	logFilename := logDir + "/logs_" + time.Now().Format("2006-01-02") + ".log"
	logFile, err := os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		log.SetOutput(logFile)
	}

	log.Println(message)
}

func CalcAbsoluteGainByCycle(cycle *database.Cycle) float64 {

	quantity := cycle.Quantity
	buyPrice := cycle.Buy.Price
	sellPrice := cycle.Sell.Price

	buyTotal := quantity * buyPrice
	sellTotal := quantity * sellPrice
	gain := sellTotal - buyTotal

	return gain
}

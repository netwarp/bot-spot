package mexc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	APIKey    string
	APISecret string
	BaseURL   string
}

func NewClient() *Client {
	return &Client{
		APIKey:    os.Getenv("MEXC_API_KEY"),
		APISecret: os.Getenv("MEXC_SECRET_KEY"),
		BaseURL:   "https://api.mexc.com",
	}
}

func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}

// Generates HMAC SHA256 signature for a signed request
func (c *Client) signRequest(queryString string) string {
	h := hmac.New(sha256.New, []byte(c.APISecret))
	h.Write([]byte(queryString))
	return hex.EncodeToString(h.Sum(nil))
}

// Sends an HTTP request and returns the response body
func (c *Client) sendRequest(method, endpoint, queryString string) ([]byte, error) {
	fullURL := fmt.Sprintf("%s%s?%s", c.BaseURL, endpoint, queryString)

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MEXC-APIKEY", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("Raw API Response:", string(body))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: HTTP status %d - %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *Client) CheckConnection() {
	_, err := c.sendRequest("GET", "/api/v3/ping", "")
	if err != nil {
		log.Fatalf("Failed to connect to MEXC: %v", err)
	}

	color.Green("Connected to MEXC API successfully")
	fmt.Println("")
}

func (c *Client) GetBalanceUSD() (float64, error) {
	color.Blue("Checking USDC balance...")

	timestamp := time.Now().UnixMilli()
	queryString := fmt.Sprintf("timestamp=%d", timestamp)
	signature := c.signRequest(queryString)
	signedQuery := fmt.Sprintf("%s&signature=%s", queryString, signature)

	body, err := c.sendRequest("GET", "/api/v3/account", signedQuery)
	if err != nil {
		log.Fatalf("Error fetching balance: %v", err)
	}

	balances, _, _, err := jsonparser.Get(body, "balances")
	if err != nil {
		log.Fatalf("Error getting balances: %v", err)
	}

	var freeFloat float64
	_, err = jsonparser.ArrayEach(balances, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		asset, _ := jsonparser.GetString(value, "asset")
		if asset == "USDC" {
			freeStr, _ := jsonparser.GetString(value, "free")
			free, _ := strconv.ParseFloat(freeStr, 64)
			freeFloat = free
		}
	})

	return freeFloat, nil
}

func (c *Client) GetLastPriceBTC() (float64, error) {
	queryString := "symbol=BTCUSDC"
	body, err := c.sendRequest("GET", "/api/v3/ticker/price", queryString)
	if err != nil {
		log.Fatalf("Error fetching BTC price: %v", err)
	}

	priceStr, err := jsonparser.GetString(body, "price")
	if err != nil {
		log.Fatalf("Error extracting price: %v", err)
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Fatalf("Error converting price: %v", err)
	}

	return price, nil
}

func (c *Client) CreateOrder(side string, price, quantity string) ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	queryString := fmt.Sprintf(
		"symbol=BTCUSDC&side=%s&type=LIMIT&quantity=%s&price=%s&timestamp=%s",
		side, quantity, price, timestamp,
	)

	signature := c.signRequest(queryString)
	signedQuery := fmt.Sprintf("%s&signature=%s", queryString, signature)

	// Send request
	body, err := c.sendRequest("POST", "/api/v3/order", signedQuery)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	return body, nil
}

func (c *Client) GetOrderById(id string) ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	queryString := fmt.Sprintf("symbol=BTCUSDC&orderId=%s&timestamp=%s", id, timestamp)
	signature := c.signRequest(queryString)
	signedQuery := fmt.Sprintf("%s&signature=%s", queryString, signature)

	// Send request
	body, err := c.sendRequest("GET", "/api/v3/order", signedQuery)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	//fmt.Println("Raw API Response:", string(body))

	return body, nil
}

func (c *Client) IsFilled(order string) (bool, error) {
	status, err := jsonparser.GetString([]byte(order), "status")
	if err != nil {
		return false, fmt.Errorf("failed to parse order status: %w", err)
	}

	return status == "FILLED", nil
}

func (c *Client) CancelOrder(orderID string) ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	queryString := fmt.Sprintf("symbol=BTCUSDC&orderId=%s&timestamp=%s", orderID, timestamp)
	signature := c.signRequest(queryString)
	signedQuery := fmt.Sprintf("%s&signature=%s", queryString, signature)

	body, err := c.sendRequest("DELETE", "/api/v3/order", signedQuery)
	if err != nil {
		return nil, fmt.Errorf("error canceling order %s: %v", orderID, err)
	}

	color.Green("Order %s canceled successfully", orderID)
	return body, nil
}

func (c *Client) GetOpenOrders() ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	queryString := fmt.Sprintf("symbol=BTCUSDC&timestamp=%s", timestamp)
	signature := c.signRequest(queryString)
	signedQuery := fmt.Sprintf("%s&signature=%s", queryString, signature)

	body, err := c.sendRequest("GET", "/api/v3/openOrders", signedQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching open orders: %v", err)
	}

	return body, nil
}

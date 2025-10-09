package tools

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"testing"
)

const ConfigFilename = "../bot.conf"

func TestTelegram(t *testing.T) {
	err := godotenv.Load(ConfigFilename)
	if err != nil {
		log.Fatal("Error loading config file")
	}

	var message = ""
	message += fmt.Sprintf("✅ Cycle 12 completed \n\n")
	message += fmt.Sprintf("📉 Buy Price: 98000 \n\n")
	message += fmt.Sprintf("📈 Sell Price: 102000 \n\n")
	message += fmt.Sprintf("💰 Gain: 4000$ \n\n")
	Telegram(message)
}

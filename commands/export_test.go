package commands

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestExport(t *testing.T) {
	_ = godotenv.Load("../bot.conf")

	Export()
}

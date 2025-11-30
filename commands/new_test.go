package commands

import (
	"fmt"
	"testing"
)

func TestCalcAmountUSD(t *testing.T) {

	amountUSD := CalcAmountUSD(200.32, 6.0)
	fmt.Println(amountUSD)

	priceBTC := 98000.00
	availableUSD := 10.00

	amountCycleBTC := CalcAmountBTC(availableUSD, priceBTC)
	fmt.Println(amountCycleBTC)
}

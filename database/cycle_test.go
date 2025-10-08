package database

import (
	"log"
	"testing"
)

func TestCycleNew(t *testing.T) {
	cycle := &Cycle{
		Exchange: "mexc",
		Status:   Buy,
		Quantity: 100,
		Buy: BuyStruct{
			Price: 98000,
			ID:    "123456789",
		},
		Sell: SellStruct{
			Price: 102000,
			ID:    "123456789",
		},
	}

	newCycle, err := CycleNew(cycle)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(newCycle)
}

func TestCycleList(t *testing.T) {
	cycles, err := CycleList()
	if err != nil {
		t.Fatal(err)
	}
	for _, cycle := range cycles {
		log.Println(cycle)
	}
}

func TestCycleGetById(t *testing.T) {
	id := 1
	cycle, err := CycleGetById(id)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(cycle)
}

func TestCycleDeleteById(t *testing.T) {
	id := 1
	err := CycleDeleteById(id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCycleListPerPage(t *testing.T) {
	page := 1
	perPage := 10
	cycles, err := CycleListPerPage(page, perPage)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cycles)
}

func TestCycleUpdate(t *testing.T) {
	id := 1
	field := "status"
	value := Sell

	result, err := CycleUpdate(id, field, value)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

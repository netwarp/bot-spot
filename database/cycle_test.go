package database_test

import (
	"github.com/joho/godotenv"
	"log"
	"main/commands"
	"main/database"
	"testing"
)

func TestCycleNew(t *testing.T) {
	godotenv.Load("../bot.conf")

	cycle, err := commands.PrepareNewCycle()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", cycle)

	cycle.Buy.ID = "xxx"
	cycle.Status = database.Buy

	id, err := database.CycleNew(cycle)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
}

func TestCycleList(t *testing.T) {
	cycles, err := database.CycleList()
	if err != nil {
		t.Fatal(err)
	}
	for _, cycle := range cycles {
		log.Println(cycle)
	}
}

func TestCycleGetById(t *testing.T) {
	id := 1
	cycle, err := database.CycleGetById(id)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(cycle)
}

func TestCycleDeleteById(t *testing.T) {
	id := 2
	err := database.CycleDeleteById(id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCycleListPerPage(t *testing.T) {
	page := 1
	perPage := 10
	cycles, err := database.CycleListPerPage(page, perPage)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cycles)
}

func TestCycleUpdate(t *testing.T) {
	id := 1
	field := "status"
	value := database.Sell

	result, err := database.CycleUpdate(id, field, value)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

package main

import (
	"errors"
	"fmt"
	"log"
	"main/commands"
	"main/database"
	"os"
)

const version = "v3.1.0"

func menu() {
	fmt.Println("")
	fmt.Println("Cryptomancien - BOT SPOT - " + version + " - beta")
	fmt.Println("")
	fmt.Println("--new			-n		Start new cycle")
	fmt.Println("--update		-u		Update running cycles")
	fmt.Println("--server		-s		Start local server")
	fmt.Println("--cancel		-c		Cancel cycle by id - Example: -c 123")
	fmt.Println("--auto			-a		Mode auto")
	fmt.Println("--clear 		-cl		Clear range (start end) - Example: -cl 12 36")
	fmt.Println("--export		-e		Export CSV file")
	fmt.Println("--restore		-r		Restore database from JSON file")
	fmt.Println("")
}

func initialize() error {
	commands.CreateConfigFileIfNotExists()
	commands.LoadDotEnv()

	// Create an exports folder if not exists
	path := "exports"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}

	err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	err := initialize()
	if err != nil {
		panic("Error initializing database: %v")
	}

	args := os.Args[1:]
	if len(args) == 0 {
		menu()
		return
	}

	cmd := args[0]
	switch cmd {
	case "--new", "-n":
		err := commands.New()
		if err != nil {
			log.Fatal(err)
		}
		break
	case "--update", "-u":
		err := commands.Update()
		if err != nil {
			log.Fatal(err)
		}
		break
	case "--server", "-s":
		err := commands.Server()
		if err != nil {
			log.Fatal(err)
		}
		break
	case "--cancel", "-c":
		err := commands.Cancel()
		if err != nil {
			log.Fatal(err)
		}
		break
	case "--clear", "-cl":
		err := commands.Clear()
		if err != nil {
			log.Fatal(err)
		}
	case "--auto", "-a":
		commands.Auto()
		break
	case "--export", "-e":
		commands.Export(true)
		break
	default:
		menu()
		return
	}

	//	"--export": func() { commands.Export(true) },
	//	"-e":       func() { commands.Export(true) },

}

package main

import (
	"errors"
	"fmt"
	"log"
	"main/commands"
	"main/database"
	"os"
	"slices"
)

const version = "v3.0.8"

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

	_, err := database.InitDatabase()
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

	actions := map[string]func(){
		"--new":     commands.New,
		"-n":        commands.New,
		"--update":  commands.Update,
		"-u":        commands.Update,
		"--server":  commands.Server,
		"-s":        commands.Server,
		"--cancel":  commands.Cancel,
		"-c":        commands.Cancel,
		"--auto":    commands.Auto,
		"-a":        commands.Auto,
		"--clear":   commands.Clear,
		"-cl":       commands.Clear,
		"--export":  func() { commands.Export(true) },
		"-e":        func() { commands.Export(true) },
		"--restore": commands.Restore,
		"-r":        commands.Restore,
	}

	for key, action := range actions {
		if slices.Contains(args, key) {
			action()
			return
		}
	}

	menu()
}

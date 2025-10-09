package commands

import (
	"fmt"
	"github.com/fatih/color"
	"main/database"
	"os"
	"strconv"
)

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func Clear() error {
	args := os.Args[2:]

	if len(args) != 2 {
		color.Red("Start and End required")
		color.Cyan("Example: go run . -cl 12 35")
		return nil
	}

	startStr := args[0]
	endStr := args[1]

	start, _ := strconv.Atoi(startStr)
	end, _ := strconv.Atoi(endStr)

	if start == end {
		color.Yellow("Delete one cycle %d", start)

		err := database.CycleDeleteById(start)
		if err != nil {
			return fmt.Errorf("error deleting cycle %d: %v", start, err)
		}

		color.Green("Cycle %d successfully deleted", start)
		return nil
	}

	r := makeRange(start, end)

	for i := range r {
		color.White("Deleting %d", r[i])
		//database.DeleteByIdInt(int32(r[i]))
		err := database.CycleDeleteById(r[i])
		if err != nil {
			return fmt.Errorf("error deleting cycle %d: %v", r[i], err)
		}
	}

	color.Green("Range successfully deleted")
	return nil
}

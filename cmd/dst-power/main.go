package main

import (
	"fmt"
	"os"

	"github.com/MainPoser/dst-power/cmd/dst-power/app"
)

func main() {
	command := app.NewPowerCommand()
	err := command.Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

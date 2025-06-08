package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/app"
	"github.com/madalinpopa/gocost/internal/config"
	"github.com/madalinpopa/gocost/internal/data"
	"github.com/spf13/viper"
)

func main() {

	if err := config.LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config file: %v", err)
		os.Exit(1)
	}

	dataFilePath := viper.GetString(config.DataFileField)
	defaultCurrency := viper.GetString(config.CurrencyField)
	initialData, err := data.LoadData(dataFilePath, defaultCurrency)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading data from %s: %v", dataFilePath, err)
		os.Exit(1)
	}

	app := app.New(initialData, dataFilePath)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

}

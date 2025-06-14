package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madalinpopa/gocost/internal/app"
	"github.com/madalinpopa/gocost/internal/config"
	"github.com/madalinpopa/gocost/internal/data"
	"github.com/spf13/viper"
)

// version will be set during build time
var version = "dev"

func main() {
	// Define version flag
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	// If version flag is set, print version and exit
	if *versionFlag {
		fmt.Printf("gocost version %s\n", version)
		os.Exit(0)
	}

	if err := config.LoadConfig(); err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "Error loading config file: %v", err); err != nil {
			os.Exit(2)
		}
		os.Exit(1)
	}

	dataFilePath := viper.GetString(config.DataFileField)
	defaultCurrency := viper.GetString(config.CurrencyField)
	initialData, err := data.LoadData(dataFilePath, defaultCurrency)
	if err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "Error loading data from %s: %v", dataFilePath, err); err != nil {
			os.Exit(2)
		}
		os.Exit(1)
	}

	a := app.New(initialData, dataFilePath)
	p := tea.NewProgram(a, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "Error running program: %v\n", err); err != nil {
			os.Exit(2)
		}
		os.Exit(1)
	}

}

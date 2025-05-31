package ui

import (
	"time"

	"github.com/madalinpopa/gocost/internal/data"
)

type Data struct {
	Data         *data.DataRoot
	DataFilePath string
}

type MonthYear struct {
	CurrentMonth time.Month
	CurrentYear  int
}

type WindowSize struct {
	Width  int
	Weight int
}

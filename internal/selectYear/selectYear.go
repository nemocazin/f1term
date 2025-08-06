package selectYear

import (
	"github.com/pterm/pterm"
	"strconv"
)

func printYears() {
	str := ""
	for year := 2022; year <= 2023; year++ {
		str += strconv.Itoa(year) + "\n"
	}
	pterm.DefaultBox.
		WithTitle("Select your seasons").
		WithTitleBottomRight().
		Println(str)
}

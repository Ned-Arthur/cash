package main

import (
	"fmt"
	"os"
	"errors"
	"strconv"
	"net/http"
	"encoding/json"

	"github.com/charmbracelet/huh"
	"github.com/joho/godotenv"		// Read .env file for API key
)


var symbols = map[string]string {
	"AUD": "$",
	"USD": "$",
	"GBP": "£",
	"EUR": "€",
}

type oerData struct {
	disclaimer string `json:"disclaimer"`
	license string `json:"license"`
	timestamp string `json:"timestamp"`
	Base string `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

/* HELPER FUNCTIONS */

// Make a request to the Open Exchange Rates API and decode the JSON response
func makeReq() oerData {
	oer_token := os.Getenv("OER_APPID")

	resp, err := http.Get("https://openexchangerates.org/api/latest.json?app_id="+oer_token)
	if err != nil {
		panic(err)
	}

	var data oerData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	return data
}

func main() {
	godotenv.Load()

	fromCur, toCur, amt := runMenu()

	data := makeReq()

	rateFrom := data.Rates[fromCur]
	rateTo := data.Rates[toCur]

	converted := amt * (rateTo / rateFrom)

	fmt.Printf("%s%.2f %s converts to %s%.2f %s\n", symbols[fromCur], amt, fromCur, symbols[toCur], converted, toCur)
}

// Get our user input using the huh library
func runMenu() (string, string, float64) {
	var fromCur string
	var toCur string
	var amt string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("From currency:").
				Options(
					huh.NewOption("$ AUD", "AUD"),
					huh.NewOption("$ USD", "USD"),
					huh.NewOption("£ GBP", "GBP"),
					huh.NewOption("€ EUR", "EUR"),
				).
				Value(&fromCur),

			huh.NewSelect[string]().
				Title("To currency:").
				Options(
					huh.NewOption("$ AUD", "AUD"),
					huh.NewOption("$ USD", "USD"),
					huh.NewOption("£ GBP", "GBP"),
					huh.NewOption("€ EUR", "EUR"),
				).
				Value(&toCur),
			huh.NewInput().
				Title("Amount to convert").
				Value(&amt).
				Validate(func(str string) error {
					if _, err := strconv.ParseFloat(str, 64); err != nil {
						return errors.New("Please enter a number, eg. 3.50")
					}
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		panic(err)
	}

	amtNum, _ := strconv.ParseFloat(amt, 64)

	return fromCur, toCur, amtNum
}


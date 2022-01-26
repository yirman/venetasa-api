package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

const dateHourRateRegex = "🗓️ ([0-9]+(/[0-9]+)+)🕒 (0?[0-9]|1[0-9]|2[0-3]):[0-9]+ [a-zA-Z]+💵 [a-zA-Z]+\\. [0-9]+,[0-9]+"

const rateRegex = "[0-9]+,[0-9]+(🔺|🔻|=)"

type Rate struct {
	Data string `json:"data"`
}

type Rates []Rate

func main() {

	fmt.Println("Iniciando server")

	r := mux.NewRouter()
	r.HandleFunc("/rates", queryRates).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

func queryRates(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

	c := colly.NewCollector()

	c.OnHTML("section.tgme_channel_history", func(e *colly.HTMLElement) {

		var rates Rates
		e.ForEach("div.tgme_widget_message_wrap", func(_ int, f *colly.HTMLElement) {
			message := f.DOM.Find("div.tgme_widget_message_text").Text()
			match1, _ := regexp.MatchString(rateRegex, message)
			match2, _ := regexp.MatchString(dateHourRateRegex, message)
			if match1 || match2 {
				rates = append(rates, Rate{Data: message})
				fmt.Println(message + "\n" + "--------------------------------------------")
			}
		})

		jsonResponse, err := json.Marshal(rates)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	})
	c.Visit("https://t.me/s/enparalelovzlatelegram")
}

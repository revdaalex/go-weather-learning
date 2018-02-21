package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"html/template"
)

type WeatherTemp struct {
	Temp     float64 `json:"temp"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
	Humidity int32   `json:"humidity"`
	Pressure float64 `json:"pressure"`
}

type WeatherWind struct {
	Speed float64 `json:"speed"`
	Deg   float64   `json:"deg"`
} 

type WeatherDescr struct {
	Icon string `json:"icon"`
	Full string `json:"description"`
}

type WeatherBase struct {
	Base    string         `json:"base"`
	Main    WeatherTemp    `json:"main"`
	Name    string         `json:"name"`
	Wind    WeatherWind    `json:"wind"`
	Descr   []WeatherDescr `json:"weather"`
}

func main() {
	tmpl, err := template.New("template.html").ParseFiles("template.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		c := http.Client{}
		resp, err := c.Get("http://api.openweathermap.org/data/2.5/weather?q=Yekaterinburg,ru&lang=ru&units=metric&appid=b0e8c750497d3d6add4e1b144715e5b2")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		weather := WeatherBase{}

		err = json.Unmarshal(body, &weather)
		if err != nil {
			log.Fatal(err)
		}

		err = tmpl.Execute(w, weather)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "style.css")
	})

	http.ListenAndServe(":8081", nil)
}

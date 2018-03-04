package main

import (
	"html/template"
	"net/http"
	"io"
	"github.com/labstack/echo"
	"io/ioutil"
	"encoding/json"
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
	Deg   float64 `json:"deg"`
}

type WeatherDescr struct {
	Icon string `json:"icon"`
	Full string `json:"description"`
}

type WeatherBase struct {
	Base  string         `json:"base"`
	Main  WeatherTemp    `json:"main"`
	Name  string         `json:"name"`
	Wind  WeatherWind    `json:"wind"`
	Descr []WeatherDescr `json:"weather"`
}

type Template struct {
	templates *template.Template
}

var url string

var City string

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func indexHandler(c echo.Context) error  {
	defer c.Request().Body.Close()
	cl := http.Client{}
	City = "Yekaterinburg"

	if c.Request().Method == http.MethodPost {
		City = c.FormValue("city")
	}

	if City == "Revda" {
		url = "http://api.openweathermap.org/data/2.5/weather?id=502011&lang=ru&units=metric&appid=b0e8c750497d3d6add4e1b144715e5b2"
	} else {
		url = "http://api.openweathermap.org/data/2.5/weather?q=" + City + ",ru&lang=ru&units=metric&appid=b0e8c750497d3d6add4e1b144715e5b2"
	}

	resp, err := cl.Get(url)
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return nil
		}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	weather := WeatherBase{}

	err = json.Unmarshal(body, &weather)
	if err != nil {
		c.Logger().Fatal(err)
	}
	return c.Render(http.StatusOK, "template", weather)
}

func main()  {
	t := &Template{
		templates: template.Must(template.ParseFiles("template.html")),
	}
	e := echo.New()
	e.Renderer = t
	e.GET("/", indexHandler)
	e.POST("/City", indexHandler)
	e.Static("/assets/style.css", "assets/style.css")

	e.Logger.Fatal(e.Start(":8081"))
}

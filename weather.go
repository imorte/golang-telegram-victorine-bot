package main

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strings"
	"strconv"
)

type Content struct {
	Query struct {
		Results struct {
			Channel struct {
				Item struct {
					Condition struct {
						Temp string `json:"temp"`
					} `json:"condition"`
				} `json:"item"`
			} `json:"channel"`
		} `json:"results"`
	} `json:"query"`
}

func getWeather() {
	var fahr Content
	var fahrOy Content
	weatherMoscow := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20weather.forecast%20where%20woeid%20in%20(2122265)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"
	weatherOymyakon := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20weather.forecast%20where%20woeid%20in%20(select%20woeid%20from%20geo.places(1)%20where%20text%3D%22oymyakon%22)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"

	res, err := http.Get(weatherMoscow)

	if err != nil {
		err.Error()
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		err.Error()
	}
	_ = json.Unmarshal(body, &fahr)

	res, err = http.Get(weatherOymyakon)

	if err != nil {
		err.Error()
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		err.Error()
	}
	_ = json.Unmarshal(body, &fahrOy)

	var result int

	// to int
	result := fahr.Query.Results.Channel.Item.Condition.Temp + fahrOy.Query.Results.Channel.Item.Condition.Temp) / 2

	println(result)

}

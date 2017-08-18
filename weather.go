package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

func (c Content) Convert() (converted int) {
	str := c.Query.Results.Channel.Item.Condition.Temp
	var err error
	converted, err = strconv.Atoi(str)
	if err != nil {
		log.Println(err)
	}
	return
}

func getWeather() {
	var fahr Content
	var fahrOy Content
	weatherMoscow := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20weather.forecast%20where%20woeid%20in%20(2122265)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"
	weatherOymyakon := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20weather.forecast%20where%20woeid%20in%20(select%20woeid%20from%20geo.places(1)%20where%20text%3D%22oymyakon%22)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"

	res, err := http.Get(weatherMoscow)

	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(body, &fahr)
	if err != nil {
		log.Fatal(err)
	}

	res, err = http.Get(weatherOymyakon)

	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(body, &fahrOy)
	if err != nil {
		log.Fatal(err)
	}

	var result int
	result = (fahr.Convert() + fahrOy.Convert()) / 2

	println(result)

}

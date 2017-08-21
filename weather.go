package main

import (
	"github.com/tidwall/gjson"

	"io/ioutil"
	"log"
	"net/http"
)

func getWeather(int, int) {
	weatherMoscow := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20weather.forecast%20where%20woeid%20in%20(2122265)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"
	weatherOymyakon := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20weather.forecast%20where%20woeid%20in%20(select%20woeid%20from%20geo.places(1)%20where%20text%3D%22oymyakon%22)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"

	res, err := http.Get(weatherMoscow)

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	fahr := gjson.GetBytes(body, "query.results.channel.item.condition.temp")

	res, err = http.Get(weatherOymyakon)

	if err != nil {
		panic(err)
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	fahrOy := gjson.GetBytes(body, "query.results.channel.item.condition.temp")

	return int(fahr.Int()), int(fahrOy.Int())
}

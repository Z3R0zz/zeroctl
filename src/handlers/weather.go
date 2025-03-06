package handlers

import (
	"errors"
	"fmt"
	"os"
	"zeroctl/src/database"
	"zeroctl/src/utils"

	"github.com/valyala/fasthttp"
)

func GetWeather() (*fasthttp.Response, error) {
	if os.Getenv("OPENWEATHER_API_KEY") == "" || os.Getenv("OPENWEATHER_CITY_ID") == "" || os.Getenv("OPENWEATHER_UNITS") == "" {
		return nil, errors.New("OPENWEATHER_API_KEY, OPENWEATHER_CITY_ID, and OPENWEATHER_UNITS must be set")
	}

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?APPID=%s&id=%s&units=%s", os.Getenv("OPENWEATHER_API_KEY"), os.Getenv("OPENWEATHER_CITY_ID"), os.Getenv("OPENWEATHER_UNITS"))
	resp, err := utils.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather: %w", err)
	}

	return resp, nil
}

func CacheWeatherData() error {
	data, err := GetWeather()
	if err != nil {
		return err
	}

	err = database.StoreJsonData("weather", data)
	if err != nil {
		return err
	}

	return nil
}

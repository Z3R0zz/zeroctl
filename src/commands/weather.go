package commands

import (
	"encoding/json"
	"zeroctl/src/handlers"
	"zeroctl/src/types"
)

func init() {
	types.RegisterCommand(types.Command{
		Name:        "weather",
		Description: "Get the current weather",
		Handler: func() string {
			weather, err := handlers.GetWeather()
			if err != nil {
				return "Failed to get weather: " + err.Error() + "\n"
			}

			weatherString, err := json.Marshal(weather)
			if err != nil {
				return "Failed to marshal weather data: " + err.Error() + "\n"
			}

			return string(weatherString)
		},
	})
}

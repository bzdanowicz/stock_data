package internal

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	ApiKey       string   `json:"apiKey"`
	UserQuotes   []string `json:"quotes"`
	BaseCurrency string   `json:"base"`
}

func ReadConfiguration() Configuration {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	return configuration
}

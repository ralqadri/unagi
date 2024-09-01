package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Token     string `json:"token"`
	BotPrefix string `json:"botPrefix"`
}

// TODO: store these in env instead
func ReadConfig() (*Config, error) {
	log.Println("Reading 'config.json'...")
	data, err := os.ReadFile("config/config.json")
	if err != nil {
		log.Fatalf("Error reading config file!: %v", err)
	}

	log.Println("Unmarshaling 'config.json' ...")

	var cfg Config
	err = json.Unmarshal([]byte(data), &cfg) // slicing the data into a byte slice then unmarshaling it into the cfg variable
	if err != nil {
		log.Printf("Error unmarshaling config file!: %v", err)
	}
	log.Printf("Config file succesfully read!\n")

	return &cfg, nil
}

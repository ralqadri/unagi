package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Token		string	`json:"token"`
	BotPrefix	string	`json:"botPrefix"`
}

func ReadConfig() (*Config, error) {
	fmt.Println("reading config.json ...")
	data, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("error reading config.json!: %v", err)
		return nil, err
	}
	fmt.Println("unmarshaling config.json ...")
	var cfg Config
	err = json.Unmarshal([]byte(data), &cfg) // slicing the data into a byte slice then unmarshaling it into the cfg variable
	if err != nil {
		fmt.Printf("error unmarshaling config.json!: %v", err)
		return nil, err
	}
	return &cfg, nil
}
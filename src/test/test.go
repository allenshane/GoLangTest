package main

import (
	"encoding/json"
	"fmt"
	"github/copy"
	"os"
)

type Config struct {
	Folders struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	} `json:"folders"`
}

func LoadConfiguration(file string) (config Config, err error) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return
	}
	dec := json.NewDecoder(configFile)
	err = dec.Decode(&config)
	return
}

func main() {
	fmt.Println("Starting the app....")
	config, _ := LoadConfiguration(os.Args[0])
	copy.CopyDirectory(config.Folders.Source, config.Folders.Destination)
}

package ulab

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

type Config struct {
	Root string `json:"root"`
	Data string `json:"data"`
}

func ReadConfigFile() *Config {
	var filename string
	switch runtime.GOOS {
	case "linux", "darwin", "illumos":
		filename = "/opt/ulab/ulab.json"
	case "netbsd", "freebsd", "openbsd", "dragonfly":
		filename = "/usr/local/ulab/ulab.json"
	}
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
	var c *Config
	json.Unmarshal(jsonData, &c)
	return (c)
}

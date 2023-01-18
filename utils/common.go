package utils

import (
	"encoding/json"
	"os"
)

func ReadJsonFile(confPath string, conf interface{}) error {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, conf); err != nil {
		return err
	}
	return nil
}

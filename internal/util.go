package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func WriteJSON(file string, data []byte) error {
	err := ioutil.WriteFile(file, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadJSON(file string, v interface{}) error {
	data, err := os.Open(file)
	if err != nil {
		return err
	}
	byteValue, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteValue, &v)
	return nil
}

package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func WriteJSON(file string, data []byte) error {
	json.Marshal(data)
	err := ioutil.WriteFile(file, data, 0644)
	if err != nil {
		return err
	}
	log.Printf("Success Writing json file: %v", err)
	return nil
}

func ReadJSON(file string, v interface{}) error {
	data, err := os.Open(file)
	if err != nil {
		log.Printf("Fail to Open json file: %v", err)
		return err
	}
	byteValue, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteValue, &v)
	return nil
}

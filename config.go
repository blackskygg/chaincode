package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config map[string]string

func FromFile(file string) (Config, error) {
	var config Config
	dat, err := ioutil.ReadFile(file)
	err = json.Unmarshal(dat, &config)
	return config, err
}

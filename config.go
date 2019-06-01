package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type config struct {
	Folder string `json:"folder"`
}

func loadConfig(name string) (config, error) {
	home := os.Getenv("HOME")
	if _, err := os.Stat(home + "/.config/todo"); os.IsNotExist(err) {
		err := os.MkdirAll(home+"/.config/todo", os.ModePerm)
		if err != nil {
			return config{}, err
		}
	}
	j, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return config{}, err
	}

	defer j.Close()

	// Default Config
	c := config{
		Folder: home + "/.config/todo",
	}

	jb, err := ioutil.ReadAll(j)
	if err != nil {
		return config{}, err
	}
	json.Unmarshal(jb, &c)
	return c, nil
}

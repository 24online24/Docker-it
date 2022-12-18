package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// TODO check for the file, if not available, create it and fill it with default values
// validate null input as well

func get_settings() {
	dat, err := os.ReadFile(".settings")
	if err != nil {
		log.Fatal(err)
	}
	str := strings.Split(string(dat), "\n")
	refresh_rate, err = strconv.Atoi(strings.Trim(str[0], "\r"))
	if err != nil {
		log.Fatal(err)
	}
	terminal_setting = strings.Trim(str[1], "\r")
	theme_color = strings.Trim(str[2], "\r")
	if env == "windows" {
		if len(str) == 4 {
			docker_path = strings.Trim(str[3], "\r")
		} else {
			docker_path = "C:\\Program Files\\Docker\\Docker\\Docker Desktop.exe"
		}
	}
	fmt.Println("Settings have been imported succesfully!")
}

func save_settings() {
	val := fmt.Sprint(refresh_rate) + "\n" + terminal_setting + "\n" + theme_color + "\n"
	if env == "windows" {
		if docker_path == "" {
			docker_path = "C:\\Program Files\\Docker\\Docker\\Docker Desktop.exe"
		}
		val += docker_path
	}
	data := []byte(val)

	err := ioutil.WriteFile(".settings", data, 0)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Settings have been saved succesfully!")
}

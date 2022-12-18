package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

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
	if env == "windows" {
		docker_path = strings.Trim(str[2], "\r")
	}
	fmt.Println("Settings have been imported succesfully!")
}

func save_settings() {
	val := fmt.Sprint(refresh_rate) + "\n" + terminal_setting + "\n"
	if env == "windows" {
		val += docker_path
	}
	data := []byte(val)

	err := ioutil.WriteFile(".settings", data, 0)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Settings have been saved succesfully!")
}

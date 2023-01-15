package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// get settings from .settings file
func get_settings() {
	dat, err := os.ReadFile(".settings")
	if err != nil {
		fil, err2 := os.Create(".settings")
		value := fmt.Sprint(refresh_rate) + "\n" + terminal_setting + "\n" + theme_color + "\n"
		if env == "windows" {
			value += docker_path
		}
		data_b := []byte(value)
		fil.Write(data_b)
		if err2 != nil {
			log.Fatal(err2)
		}
		dat, _ = os.ReadFile(".settings")
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

// save settings to .settings file
func save_settings() {
	val := fmt.Sprint(refresh_rate) + "\n" + terminal_setting + "\n" + theme_color + "\n"
	if env == "windows" {
		val += docker_path
	}
	data := []byte(val)

	err := os.WriteFile(".settings", data, 0)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Settings have been saved succesfully!")
}

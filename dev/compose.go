package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var output string
	var space string = "  "
	file, err := os.Create("./docker-compose.yml")
	check(err)

	output = "version: 3.9\n"
	_, err = file.WriteString(output)
	check(err)

	var input string
	fmt.Println("How many services do you want?")
	fmt.Scanln(&input)

	nr_services, err := strconv.Atoi(input)
	check(err)

	for i := 1; i <= nr_services; i++ {
		fmt.Printf("What name do you want to use for service %d?\n", i)
		fmt.Scanln(&input)
		output = space + input + ":\n"
		_, err = file.WriteString(output)
		check(err)

		fmt.Printf("What image do you want to use for service %d?\n", i)
		fmt.Scanln(&input)
		output = strings.Repeat(space, 2) + "image: " + input + "\n"
		_, err = file.WriteString(output)
		check(err)

		fmt.Printf("What internal (container) port do you want to use for service %d?\nLeave empty if you don't want to expose (any more) ports.\n", i)

		input = ""
		fmt.Scanln(&input)
		if input != "" {
			output = strings.Repeat(space, 2) + "ports:\n"
			_, err = file.WriteString(output)
			check(err)
			output = strings.Repeat(space, 3) + "- \"" + input + ":"
			fmt.Printf("What external (host computer) port do you want to use?\n")
			fmt.Scanln(&input)
			output = output + input + "\""
			_, err = file.WriteString(output)
			check(err)
		}
		// fmt.Printf("What internal (container) port do you want to use for service %d?\n", i)
		// fmt.Scanln(&input)
		// output = "\t\tports: " + input + "\n"
		// _, err = file.WriteString(output)
		// check(err)
	}

	defer file.Close()
}

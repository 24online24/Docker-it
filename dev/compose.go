package main

import (
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var output string
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
		output = "\t" + input + ":\n"
		_, err = file.WriteString(output)
		check(err)

		fmt.Printf("What image do you want to use for service %d?\n", i)
		fmt.Scanln(&input)
		output = "\t\tmage: " + input + "\n"
		_, err = file.WriteString(output)
		check(err)

		// fmt.Printf("What image do you want to use for service %d?\n", i)
		// fmt.Scanln(&input)
		// output = "\t\tports: " + input + "\n"
		// _, err = file.WriteString(output)
		// check(err)
	}

	defer file.Close()
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var space string = "    "

func main() {

	outputFile := createFile("docker-compose.yml")
	defer outputFile.Close()

	w := bufio.NewWriter(outputFile)

	writeServices(outputFile, w)

}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createFile(fileName string) *os.File {
	outputFile, err := os.Create(fileName)
	handleError(err)
	return outputFile
}

func writeServices(outputFile *os.File, w *bufio.Writer) {
	writeLine(outputFile, w, "services:\n")

	fmt.Println("How many services do you want?")
	nrServices := -1
	_, err := fmt.Scanln(&nrServices)
	handleError(err)

	for i := 0; i < nrServices; i++ {
		service(outputFile, w, i)
	}
}

func service(outputFile *os.File, w *bufio.Writer, nrOfService int) {
	fmt.Printf("========================%d========================\n", nrOfService+1)

	input := getInput("What name do you want to use for this service?")
	writeLine(outputFile, w, space+input+":\n")

	container(outputFile, w)
	ports(outputFile, w)
}

func writeLine(outputFile *os.File, w *bufio.Writer, line string) {
	_, err := w.WriteString(line)
	handleError(err)
	w.Flush()
}

func getInput(message string) string {
	if message != "" {
		fmt.Println(message)
	}
	input := ""
	_, err := fmt.Scanln(&input)
	handleError(err)
	return input
}

func container(outputFile *os.File, w *bufio.Writer) {
	input := getInput("Do you want a standard image or do you have a dockerfile for this service?")
	output := ""

	switch input {
	case "image":
		input = getInput("What image do you want?")
		output = strings.Repeat(space, 2) + "image: " + input + "\n"
	case "dockerfile":
		fmt.Println("Path to the dockerfile.")
		input = ""
		_, err := fmt.Scanln(&input)
		handleError(err)
		output = strings.Repeat(space, 2) + "image: " + input + "\n"
	}
	writeLine(outputFile, w, output)
}

func ports(outputFile *os.File, w *bufio.Writer) {
	input := getInput("Do you want a have open ports on this container?")

	if input == "yes" {
		output := strings.Repeat(space, 2) + "ports:\n"
		writeLine(outputFile, w, output)

		output = strings.Repeat(space, 3) + "- '"
		input = getInput("Host machine port.")
		output += input + ":"
		input = getInput("Container port.")
		output += input
		output += "'"
		writeLine(outputFile, w, output)
	}
}

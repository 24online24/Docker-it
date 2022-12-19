package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var space string = "    "

func generateCompose(serviceName []string, imageOrFile []string, nameOrPath []string, bindPorts []bool, hostPort []string, containerPort []string) {

	outputFile := createFile("docker-compose.yml")
	defer outputFile.Close()

	w := bufio.NewWriter(outputFile)

	writeServices(outputFile, w, serviceName, imageOrFile, nameOrPath, bindPorts, hostPort, containerPort)

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

func writeServices(outputFile *os.File, w *bufio.Writer, serviceName []string, imageOrFile []string, nameOrPath []string, bindPorts []bool, hostPort []string, containerPort []string) {
	writeLine(outputFile, w, "services:\n")

	nrServices := len(serviceName)

	for i := 0; i < nrServices; i++ {
		service(outputFile, w, serviceName, imageOrFile, nameOrPath, bindPorts, hostPort, containerPort, i)
	}
}

func service(outputFile *os.File, w *bufio.Writer, serviceName []string, imageOrFile []string, nameOrPath []string, bindPorts []bool, hostPort []string, containerPort []string, index int) {

	writeLine(outputFile, w, space+serviceName[index]+":\n")
	writeContainer(outputFile, w, serviceName, imageOrFile, nameOrPath, bindPorts, hostPort, containerPort, index)
	ports(outputFile, w, serviceName, imageOrFile, nameOrPath, bindPorts, hostPort, containerPort, index)
}

func writeLine(outputFile *os.File, w *bufio.Writer, line string) {
	_, err := w.WriteString(line)
	handleError(err)
	w.Flush()
}

func writeContainer(outputFile *os.File, w *bufio.Writer, serviceName []string, imageOrFile []string, nameOrPath []string, bindPorts []bool, hostPort []string, containerPort []string, index int) {
	output := ""

	switch imageOrFile[index] {
	case "Image":
		output = strings.Repeat(space, 2) + "image: " + nameOrPath[index] + "\n"
	case "Custom":
		output = strings.Repeat(space, 2) + "image: " + nameOrPath[index] + "\n"
	}
	writeLine(outputFile, w, output)
}

func ports(outputFile *os.File, w *bufio.Writer, serviceName []string, imageOrFile []string, nameOrPath []string, bindPorts []bool, hostPort []string, containerPort []string, index int) {

	if bindPorts[index] {
		output := strings.Repeat(space, 2) + "ports:\n"
		writeLine(outputFile, w, output)

		output = strings.Repeat(space, 3) + "- '"
		output += hostPort[index] + ":"
		output += containerPort[index]
		output += "'\n"
		writeLine(outputFile, w, output)
	}
}

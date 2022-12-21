package main

import (
	"bufio"
	"os"
	"strings"
)

var space string = "    "

func generateCompose(cIL []containerInfo) {

	outputFile := createFile("docker-compose.yml")
	defer outputFile.Close()

	w := bufio.NewWriter(outputFile)

	writeServices(outputFile, w, cIL)

}

func createFile(fileName string) *os.File {
	outputFile, err := os.Create(fileName)
	handleError(err)
	return outputFile
}

func writeServices(outputFile *os.File, w *bufio.Writer, cIL []containerInfo) {
	writeLine(outputFile, w, "services:\n")

	nrServices := len(cIL)

	for i := 0; i < nrServices; i++ {
		service(outputFile, w, cIL[i], i)
	}
}

func service(outputFile *os.File, w *bufio.Writer, cI containerInfo, index int) {

	writeLine(outputFile, w, space+cI.serviceName+":\n")
	writeContainer(outputFile, w, cI, index)
	ports(outputFile, w, cI, index)
}

func writeLine(outputFile *os.File, w *bufio.Writer, line string) {
	_, err := w.WriteString(line)
	handleError(err)
	w.Flush()
}

func writeContainer(outputFile *os.File, w *bufio.Writer, cI containerInfo, index int) {
	output := ""

	switch cI.imageOrFile {
	case "Image":
		output = strings.Repeat(space, 2) + "image: " + cI.nameOrPath + "\n"
	case "Custom":
		output = strings.Repeat(space, 2) + "image: " + cI.nameOrPath + "\n"
	}
	writeLine(outputFile, w, output)
}

func ports(outputFile *os.File, w *bufio.Writer, cI containerInfo, index int) {

	if cI.bindPorts {
		output := strings.Repeat(space, 2) + "ports:\n"
		writeLine(outputFile, w, output)

		output = strings.Repeat(space, 3) + "- '"
		output += cI.hostPort + ":"
		output += cI.containerPort
		output += "'\n"
		writeLine(outputFile, w, output)
	}
}

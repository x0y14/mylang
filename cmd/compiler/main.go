package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 2 {
		log.Fatal("[input file path] [output file path]")
	}
	inputFilePath := os.Args[1]
	outputFilePath := os.Args[2]

	// do something
	_ = inputFilePath
	//inputFileData := compile output
	inputFileData := "exit"

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}
	_, err = outputFile.Write([]byte(inputFileData))
	if err != nil {
		log.Fatalf("faied to write output file: %s", err)
	}
}

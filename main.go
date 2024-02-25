package main

import (
	"log"
	"os"

	"github.com/LawrenceVeigas/gosqlfmt2/gosqlfmt2"
)

func main() {
	fileName := os.Args[1]
	finalQuery := gosqlfmt2.Parse(fileName)

	err := os.WriteFile(fileName, []byte(finalQuery), 0644)
	if err != nil {
		log.Fatalf("Error writing to file..\n%v", err)
	}
}

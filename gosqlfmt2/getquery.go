package gosqlfmt2

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func GetQuery(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening %s\n%s", fileName, err)
	}

	reader := bufio.NewReader(file)

	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("Error reading %s\n%s", fileName, err)
	}

	contents := string(data)

	return contents
}

package main

import (
	"log"
	"os"
)

func main() {
	file := readFileString("C:/Users/bartek/go/src/data/polish.txt")
	_, _ = os.Stdout.WriteString(file)

}


// Read entire file's contents and return a string
func readFileString(fp string) string {
	file, err := os.ReadFile(fp)	// ioutil.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	return string(file)
}

// Read entire file's contents and return bytes
func readFileBytes(fp string) []byte {
	file, err := os.ReadFile(fp)	// ioutil.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	return file		// fmt.Printf("%s", file)
}


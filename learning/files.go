// working with files
package main

import (
	"fmt"
	"log"

	//"fmt"
	"os"
	//"path/filepath"
)


const dataDir = "C:/Users/bartek/go/src/gosrc/learning/data"


func main() {
	//fmt.Println(openFile(filepath.Join(dataDir, "polish.txt")))
	//fp := filepath.Join(dataDir, "polish.txt")
	//createFileSafeWrite(fp, "write text")

}

// check file exists
func fileExists(fp string) bool {
	_, err := os.Stat(fp)
	if os.IsNotExist(err) {
		return false	// not found
	}
	return true		// file found
}


// create file if not exist
func createFileSafeWrite(fp string, text string) {
	if !fileExists(fp) {
		file, err := os.Create(fp)
		if err != nil {
			log.Panic(err)
		}
		defer file.Close()
		size, _ := file.WriteString(text)
		fmt.Printf("Wrote %v bytes", size)
	}
}


//func readFileBytes(fp string) []byte {
//	text, _ := ioutil.ReadFile(fp)
//	return text		// fmt.Printf("%s", text)
//}



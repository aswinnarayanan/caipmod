package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// err := godotenv.Load("global.start")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// DICOM_SERVER_NAME_2 := os.Getenv("DICOM_SERVER_NAME_2")
	// fmt.Println(DICOM_SERVER_NAME_2)

	//////////////////////////////////////////////////////////////

	file, err := os.Open("global.start")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var DATA_SOURCE_NOT_ACTIVE = make(map[int]string)
	var DATA_BASE_CONNECTION = make(map[int]string)
	var DATA_BASE_PATH = make(map[int]string)

	var i int
	var val string

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {

		i, val = extractVal(scanner.Text(), "DATA_SOURCE_NOT_ACTIVE_")
		if i >= 0 {
			DATA_SOURCE_NOT_ACTIVE[i] = val
		}
		i, val = extractVal(scanner.Text(), "DATA_BASE_CONNECTION_")
		if i >= 0 {
			DATA_BASE_CONNECTION[i] = val
		}
		i, val = extractVal(scanner.Text(), "DATA_BASE_PATH_")
		if i >= 0 {
			DATA_BASE_PATH[i] = val
		}
	}

	fmt.Println("Following Databases found")
	for i, _ = range DATA_SOURCE_NOT_ACTIVE {
		val, _ = strings.CutPrefix(DATA_BASE_CONNECTION[i], "jdbc\\:derby\\:")
		fmt.Printf("[%d]\n%s\n%s\n\n", i, val, DATA_BASE_PATH[i])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func extractVal(line string, prefix string) (i int, val string) {
	var s string
	i = -1
	val = ""
	if strings.HasPrefix(line, prefix) {
		tmp, _ := strings.CutPrefix(line, prefix)
		s, val, _ = strings.Cut(tmp, "=")
		// fmt.Println(s, ":", val)
		i, _ = strconv.Atoi(s)
	}
	return i, val
}

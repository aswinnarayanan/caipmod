package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {
	e, _ := os.Executable()
	fmt.Println(path.Dir(e))

	file, err := os.Open(path.Dir(e) + "/global.start")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var DATA_SOURCE_NOT_ACTIVE = make(map[int]string)
	var DATA_BASE_CONNECTION = make(map[int]string)
	var DATA_BASE_PATH = make(map[int]string)

	var i int
	var val string
	var dbi int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		i, val = extractVal(scanner.Text(), "DATA_SOURCE_NOT_ACTIVE_")
		if i > 0 {
			DATA_SOURCE_NOT_ACTIVE[i] = val
		}
		i, val = extractVal(scanner.Text(), "DATA_BASE_CONNECTION_")
		if i > 0 {
			DATA_BASE_CONNECTION[i] = val
		}
		i, val = extractVal(scanner.Text(), "DATA_BASE_PATH_")
		if i > 0 {
			DATA_BASE_PATH[i] = val
		}
	}

	fmt.Printf("\n-------------------------\nFollowing Databases found\n-------------------------\n\n")
	for i := range DATA_SOURCE_NOT_ACTIVE {
		// val, _ = strings.CutPrefix(DATA_BASE_CONNECTION[i], "jdbc\\:derby\\:")
		fmt.Printf("[%d]\nDATA_BASE_CONNECTION = %s\nDATA_BASE_PATH = %s\n\n", i, DATA_BASE_CONNECTION[i], DATA_BASE_PATH[i])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Take input db
	fmt.Print("-------------------\nEnter your database\n> ")
	fmt.Scanf("%d", &dbi)

	fmt.Printf("\n----------------------------\nUsing the following database\n\n")
	fmt.Printf("[%d]\n%s\n%s\n\n", dbi, DATA_BASE_CONNECTION[dbi], DATA_BASE_PATH[dbi])

}

func extractVal(line string, prefix string) (i int, val string) {
	var s string
	i = 0
	val = ""
	if strings.HasPrefix(line, prefix) {
		tmp, _ := strings.CutPrefix(line, prefix)
		s, val, _ = strings.Cut(tmp, "=")
		// fmt.Println(s, ":", val)
		i, _ = strconv.Atoi(s)
	}
	return i, val
}

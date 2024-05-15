package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	currentUser, _ := user.Current()
	currentPath, _ := os.Executable()

	fmt.Println(currentUser.Username)
	fmt.Println(path.Dir(currentPath))

	var DATA_SOURCE_NOT_ACTIVE = make(map[int]string)
	var DATA_BASE_CONNECTION = make(map[int]string)
	var DATA_BASE_PATH = make(map[int]string)

	var i int
	var val string
	var dbi int
	var license int
	var licenseflag string
	var lines []string

	// Take input license
	fmt.Print("-------------------\nEnter your license\n-------------------\n")
	fmt.Print("[1] Basic\n[2] Human\n[3] Advanced\n> ")
	fmt.Scanf("%d", &license)
	
	file, err := os.Open(path.Dir(currentPath) + "/properties/global.start")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		i, val = extractVal(line, "DATA_SOURCE_NOT_ACTIVE_")
		if i > 0 {
			DATA_SOURCE_NOT_ACTIVE[i] = val
		}
		i, val = extractVal(line, "DATA_BASE_CONNECTION_")
		if i > 0 {
			DATA_BASE_CONNECTION[i] = val
		}
		i, val = extractVal(line, "DATA_BASE_PATH_")
		if i > 0 {
			DATA_BASE_PATH[i] = val
		}
	}
	file.Close()

	fmt.Printf("\n-------------------------\nFollowing databases found\n-------------------------\n\n")
	for i := range DATA_SOURCE_NOT_ACTIVE {
		// val, _ = strings.CutPrefix(DATA_BASE_CONNECTION[i], "jdbc\\:derby\\:")
		fmt.Printf("[%d]\nDATA_BASE_CONNECTION = %s\nDATA_BASE_PATH = %s\n\n", i, DATA_BASE_CONNECTION[i], DATA_BASE_PATH[i])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Take input db
	fmt.Print("-------------------\nEnter your database\n-------------------\n> ")
	fmt.Scanf("%d", &dbi)

	fmt.Printf("\n----------------------------\nUsing the following database\n----------------------------\n\n")
	fmt.Printf("[%d]\n%s\n%s\n\n", dbi, DATA_BASE_CONNECTION[dbi], DATA_BASE_PATH[dbi])

	regex, _ := regexp.Compile("/winmounts/.*/data.cai.uq.edu.au/")
	
	file, err = os.Open(path.Dir(currentPath) + "/properties/global.start")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		i, _ = extractVal(line, "DATA_SOURCE_NOT_ACTIVE_")
		if i == dbi {
			line = "DATA_SOURCE_NOT_ACTIVE_" + strconv.Itoa(i) + "=NO"
		} else if i > 0 {
			line = "DATA_SOURCE_NOT_ACTIVE_" + strconv.Itoa(i) + "=YES"
		}
		i, _ = extractVal(line, "DATA_BASE_CONNECTION_")
		if i > 0 {
			line = regex.ReplaceAllString(line, "/winmounts/"+currentUser.Username+"/data.cai.uq.edu.au/")
		}
		i, _ = extractVal(line, "DATA_BASE_PATH_")
		if i > 0 {
			line = regex.ReplaceAllString(line, "/winmounts/"+currentUser.Username+"/data.cai.uq.edu.au/")
		}

		// fmt.Println(line)
		lines = append(lines, line)
	}
	// lines = append(lines, "\n")
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	modifiedContent := strings.Join(lines, "\n")
	err = os.WriteFile(path.Dir(currentPath) + "/properties/global.start", []byte(modifiedContent), 0644)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	 
 	if license == 2 {
	  licenseflag = "-lsn[5653@10.153.130.133]"
	} else if license == 3 {
	  licenseflag = "-lsn[5654@10.153.130.133]"
	} else {
	  licenseflag = "-lsn[5652@10.153.130.133]"
	}
	  
	fmt.Printf("Running license server %s", licenseflag)
	
	cmd := exec.Command("./java/jre/bin/java", "-Xmx62000M", "-jar", "pmod.jar", licenseflag)
	cmd.Dir = path.Dir(currentPath)
	out, err := cmd.Output()
	if err != nil {
	    fmt.Printf("%s", err)
	}
	output := string(out[:])
	fmt.Println(output)
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

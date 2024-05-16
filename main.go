package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
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
	var inputDbs []int
	var license int
	var licenseFlag string

	// Take input license
	fmt.Println()
	fmt.Println("------------------")
	fmt.Println("Enter your license")
	fmt.Println("------------------")
	fmt.Print("[1] Basic\n[2] Human\n[3] Advanced\n> ")
	fmt.Scanf("%d", &license)

	configRaw, err := os.ReadFile(path.Dir(currentPath) + "/properties/global.start")
	if err != nil {
		log.Fatal(err)
	}
	configLines := strings.Split(string(configRaw), "\n")

	for _, configLine := range configLines {
		i, val = ExtractVal(configLine, "DATA_SOURCE_NOT_ACTIVE_")
		if i > 0 {
			DATA_SOURCE_NOT_ACTIVE[i] = val
		}
		i, val = ExtractVal(configLine, "DATA_BASE_CONNECTION_")
		if i > 0 {
			DATA_BASE_CONNECTION[i] = val
		}
		i, val = ExtractVal(configLine, "DATA_BASE_PATH_")
		if i > 0 {
			DATA_BASE_PATH[i] = val
		}
	}

	fmt.Println("\n-------------------------")
	fmt.Println("Following databases found")
	fmt.Println("-------------------------")

	noaccessDbs := []string{}
	accessDbs := []string{}
	for i := range DATA_SOURCE_NOT_ACTIVE {
		// val, _ = strings.CutPrefix(DATA_BASE_CONNECTION[i], "jdbc\\:derby\\:")
		// fmt.Printf("[%d]\nDATA_BASE_CONNECTION = %s\nDATA_BASE_PATH = %s\n\n", i, DATA_BASE_CONNECTION[i], DATA_BASE_PATH[i])
		// fmt.Printf("[%d] %s\n", i, DATA_BASE_PATH[i])
		if unix.Access(DATA_BASE_PATH[i], unix.R_OK) == nil {
			accessDbs = append(noaccessDbs, "["+strconv.Itoa(i)+"] "+DATA_BASE_PATH[i])
		} else {
			noaccessDbs = append(accessDbs, "["+strconv.Itoa(i)+"] "+DATA_BASE_PATH[i])
		}

	}
	fmt.Println("\nFollowing databases are accessible")
	fmt.Println(strings.Join(accessDbs, "\n"))

	fmt.Println("\nFollowing databases are not accessible")
	fmt.Println(strings.Join(noaccessDbs, "\n"))

	// Get db(s) as input(s)
	fmt.Println("\n----------------------------------------------")
	fmt.Println("Enter your database(s) or Enter to disable all")
	fmt.Println("----------------------------------------------")
	fmt.Print("> ")
	inputDbs = GetInputSlice()

	regex, _ := regexp.Compile("/winmounts/.*/data.cai.uq.edu.au/")

	fmt.Println("\n---------------------------")
	fmt.Println("Modifying global.start file")
	fmt.Println("---------------------------")
	outputLines := []string{}
	for _, configLine := range configLines {
		outputLine := configLine
		i, _ = ExtractVal(configLine, "DATA_SOURCE_NOT_ACTIVE_")
		if slices.Contains(inputDbs, i) {
			// if _, err := os.Stat(DATA_BASE_PATH[i]); os.IsNotExist(err) {
			outputLine = "DATA_SOURCE_NOT_ACTIVE_" + strconv.Itoa(i) + "=NO"
			// if unix.Access(DATA_BASE_PATH[i], unix.R_OK) == nil {
			// 	fmt.Printf("(ON)  [%d] %s\n", i, DATA_BASE_PATH[i])
			// 	outputLine = "DATA_SOURCE_NOT_ACTIVE_" + strconv.Itoa(i) + "=NO"
			// } else {
			// 	fmt.Printf("(OFF) [%d] %s (NOACCESS)\n", i, DATA_BASE_PATH[i])
			// 	outputLine = "DATA_SOURCE_NOT_ACTIVE_" + strconv.Itoa(i) + "=YES"
			// }

		} else if i > 0 {
			fmt.Printf("(OFF) [%d] %s\n", i, DATA_BASE_PATH[i])
			outputLine = "DATA_SOURCE_NOT_ACTIVE_" + strconv.Itoa(i) + "=YES"
		}

		i, _ = ExtractVal(configLine, "DATA_BASE_CONNECTION_")
		if i > 0 {
			outputLine = regex.ReplaceAllString(configLine, "/winmounts/"+currentUser.Username+"/data.cai.uq.edu.au/")
		}
		i, _ = ExtractVal(configLine, "DATA_BASE_PATH_")
		if i > 0 {
			outputLine = regex.ReplaceAllString(configLine, "/winmounts/"+currentUser.Username+"/data.cai.uq.edu.au/")
		}
		outputLines = append(outputLines, outputLine)
	}
	outputRaw := strings.Join(outputLines, "\n")
	err = os.WriteFile(path.Dir(currentPath)+"/properties/global.start", []byte(outputRaw), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Choose license server
	if license == 2 {
		licenseFlag = "-lsn[5653@10.153.130.133]"
	} else if license == 3 {
		licenseFlag = "-lsn[5654@10.153.130.133]"
	} else {
		licenseFlag = "-lsn[5652@10.153.130.133]"
	}
	fmt.Printf("\nRunning license server %s", licenseFlag)

	// Run Pmod
	cmd := exec.Command("./java/jre/bin/java", "-Xmx62000M", "-jar", "pmod.jar", licenseFlag)
	cmd.Dir = path.Dir(currentPath)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	output := string(out[:])
	fmt.Println(output)
}

func ExtractVal(line string, prefix string) (i int, val string) {
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

func Numbers(s string) []int {
	var n []int
	for _, f := range strings.Fields(s) {
		i, err := strconv.Atoi(f)
		if err == nil {
			n = append(n, i)
		}
	}
	return n
}

func GetInputSlice() []int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return Numbers(scanner.Text())
}

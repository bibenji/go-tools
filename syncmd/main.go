package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	commandsLines := getLines("commands.txt")
	itemsLines := getLines("items.txt")
	paramsLines := getLines("params.txt")

	var wg sync.WaitGroup

	wg.Add(len(itemsLines))

	for _, itemValue := range itemsLines {
		go func(itemValue string) {
			defer wg.Done()

			for _, command := range commandsLines {
				for index, paramValue := range paramsLines {
					paramKey := fmt.Sprintf("[param%s]", strconv.Itoa(index+1))
					command = strings.Replace(command, paramKey, paramValue, -1)
				}

				command = strings.Replace(command, "[item]", itemValue, -1)

				execCmd(command)

			}
		}(itemValue)
	}

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Execution took %s", elapsed)
}

func getLines(filePath string) []string {
	file, fileErr := os.Open(filePath)

	if fileErr != nil {
		log.Fatal(fileErr)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines = []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if scannerErr := scanner.Err(); scannerErr != nil {
		log.Fatal(scannerErr)
	}

	return lines
}

func execCmd(command string) {
	parts := strings.Split(command, " ")
	main := parts[0]
	parts = parts[1:]
	log.Println("Starting:", main, parts)
	cmd := exec.Command(main, parts...)

	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Println(fmt.Sprintf("Error in command %s: %s", cmd, err))
	}

	log.Println("Done:", command)
}

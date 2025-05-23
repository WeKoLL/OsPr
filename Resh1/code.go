package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func performSorting(numbers []int64) {
	length := len(numbers)
	for outer := length - 1; outer > 0; outer-- {
		isSorted := true
		for inner := 0; inner < outer; inner++ {
			if numbers[inner] > numbers[inner+1] {
				numbers[inner], numbers[inner+1] = numbers[inner+1], numbers[inner]
				isSorted = false
			}
		}
		if isSorted {
			return
		}
	}
}

func readNumbersFromFile(filename string) ([]int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	var numberList []int64
	for _, part := range strings.Fields(line) {
		value, err := strconv.ParseInt(part, 10, 64)
		if err == nil {
			numberList = append(numberList, value)
		}
	}

	return numberList, nil
}

func writeNumbersToFile(filename string, numbers []int64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var builder strings.Builder
	for _, num := range numbers {
		builder.WriteString(fmt.Sprintf("%d ", num))
	}

	_, err = file.WriteString(builder.String())
	return err
}

func executeSortingProgram() {
	numberData, err := readNumbersFromFile("input.txt")
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	performSorting(numberData)

	err = writeNumbersToFile("output.txt", numberData)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}

	fmt.Println("Sorted numbers:", numberData)
}

func main() {
	executeSortingProgram()
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

type NumberProcessor struct {
	inputPath  string
	outputPath string
}

func (np *NumberProcessor) ExtractValues() ([]int64, error) {
	file, err := os.Open(np.inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var contentBuilder strings.Builder

	for {
		chunk, err := reader.ReadBytes(' ')
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("read error: %w", err)
		}

		contentBuilder.Write(chunk)

		if err == io.EOF {
			break
		}
	}

	return np.parseDigits(contentBuilder.String())
}

func (np *NumberProcessor) parseDigits(rawData string) ([]int64, error) {
	tokens := strings.FieldsFunc(rawData, func(r rune) bool {
		return r == ' ' || r == '\n' || r == '\t' || r == '\r'
	})

	var digits []int64
	for _, token := range tokens {
		if len(token) == 0 {
			continue
		}

		value, convErr := strconv.ParseInt(token, 10, 64)
		if convErr != nil {
			continue 
		}

		digits = append(digits, value)
	}

	return digits, nil
}

func (np *NumberProcessor) StoreResults(values []int64) error {
	file, err := os.Create(np.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	buffer := bufio.NewWriter(file)
	for _, num := range values {
		_, err := fmt.Fprintf(buffer, "%d ", num)
		if err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}

	return buffer.Flush()
}

func ExecuteNumberSorting(input, output string) error {
	processor := &NumberProcessor{
		inputPath:  input,
		outputPath: output,
	}

	numbers, err := processor.ExtractValues()
	if err != nil {
		return err
	}

	sort.Slice(numbers, func(a, b int) bool {
		return numbers[a] < numbers[b]
	})

	if err := processor.StoreResults(numbers); err != nil {
		return err
	}

	fmt.Printf("Successfully processed %d numbers\n", len(numbers))
	return nil
}

func main() {
	if err := ExecuteNumberSorting("input.txt", "output.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func Binary(a, b, c float64) (string, string) {
	D := b*b - 4*a*c

	if D >= 0 {
		x1 := (-b + math.Sqrt(D)) / (2 * a)
		x2 := (-b - math.Sqrt(D)) / (2 * a)
		return strconv.FormatFloat(x1, 'f', 2, 64), 
               strconv.FormatFloat(x2, 'f', 2, 64)
	} else {
		sqrtVal := math.Sqrt(-D) / (2 * a)
		realPart := -b / (2 * a)
		sqrtStr := strconv.FormatFloat(sqrtVal, 'f', 2, 64)
		realStr1 := strconv.FormatFloat(realPart, 'f', 1, 64)
		realStr2 := strconv.FormatFloat(realPart, 'f', 2, 64)
		x1 := realStr1 + "+" + sqrtStr + "i"
		x2 := realStr2 + "-" + sqrtStr + "i"
		return x1, x2
	}
}

func main() {
	filepath := "input.txt"
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error")
		return
	}

	nums := strings.Fields(string(content))
	if len(nums) < 3 {
		fmt.Println("Invalid input")
		return
	}

	a, err1 := strconv.ParseFloat(nums[0], 64)
	b, err2 := strconv.ParseFloat(nums[1], 64)
	c, err3 := strconv.ParseFloat(nums[2], 64)
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("Parse error")
		return
	}

	x1, x2 := Binary(a, b, c)
	fmt.Println(x1, x2)

	f, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error")
		return
	}
	defer f.Close()

	_, err = f.WriteString(x1 + "\n" + x2)
	if err != nil {
		fmt.Println("Error")
	}
}
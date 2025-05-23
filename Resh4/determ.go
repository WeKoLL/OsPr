package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MatrixOperations struct {
	sourceMatrix [][]float64
}

type MatrixResult struct {
	Determinant    float64
	Trace          float64
	Transposed     [][]float64
}

func (mo *MatrixOperations) LoadDataFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	var grid [][]float64
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		textLine := strings.TrimSpace(scanner.Text())
		if textLine == "" {
			continue
		}

		values := strings.Fields(textLine)
		row := make([]float64, len(values))

		for index, strVal := range values {
			value, convErr := strconv.ParseFloat(strVal, 64)
			if convErr != nil {
				return fmt.Errorf("ошибка преобразования значения '%s': %v", strVal, convErr)
			}
			row[index] = value
		}

		grid = append(grid, row)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	mo.sourceMatrix = grid
	return nil
}

func (mo *MatrixOperations) ValidateSquareMatrix() error {
	rows := len(mo.sourceMatrix)
	if rows == 0 {
		return fmt.Errorf("матрица пуста")
	}

	for _, row := range mo.sourceMatrix {
		if len(row) != rows {
			return fmt.Errorf("матрица не квадратная")
		}
	}
	return nil
}

func (mo *MatrixOperations) ComputeMatrixProperties() MatrixResult {
	return MatrixResult{
		Determinant:    mo.calculateDeterminant(mo.sourceMatrix),
		Trace:          mo.computeTrace(),
		Transposed:     mo.transposeMatrix(),
	}
}

func (mo *MatrixOperations) calculateDeterminant(mat [][]float64) float64 {
	size := len(mat)
	switch size {
	case 1:
		return mat[0][0]
	case 2:
		return mat[0][0]*mat[1][1] - mat[0][1]*mat[1][0]
	default:
		total := 0.0
		for currentCol := 0; currentCol < size; currentCol++ {
			subMatrix := mo.createSubMatrix(mat, currentCol)
			sign := 1.0
			if currentCol%2 == 1 {
				sign = -1.0
			}
			total += sign * mat[0][currentCol] * mo.calculateDeterminant(subMatrix)
		}
		return total
	}
}

func (mo *MatrixOperations) createSubMatrix(mat [][]float64, excludeCol int) [][]float64 {
	size := len(mat)
	subMat := make([][]float64, size-1)

	for i := 1; i < size; i++ {
		subMat[i-1] = make([]float64, 0, size-1)
		for j := 0; j < size; j++ {
			if j != excludeCol {
				subMat[i-1] = append(subMat[i-1], mat[i][j])
			}
		}
	}
	return subMat
}

func (mo *MatrixOperations) computeTrace() float64 {
	sum := 0.0
	for i := 0; i < len(mo.sourceMatrix); i++ {
		sum += mo.sourceMatrix[i][i]
	}
	return sum
}

func (mo *MatrixOperations) transposeMatrix() [][]float64 {
	size := len(mo.sourceMatrix)
	transposed := make([][]float64, size)
	for i := range transposed {
		transposed[i] = make([]float64, size)
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			transposed[j][i] = mo.sourceMatrix[i][j]
		}
	}
	return transposed
}

func SaveMatrixResults(outputPath string, result MatrixResult) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "Определитель: %.2f\n", result.Determinant)
	fmt.Fprintf(writer, "След: %.2f\n", result.Trace)
	fmt.Fprintln(writer, "Транспонированная матрица:")

	for _, row := range result.Transposed {
		for _, val := range row {
			fmt.Fprintf(writer, "%.2f ", val)
		}
		fmt.Fprintln(writer)
	}

	return nil
}

func DisplayResults(result MatrixResult) {
	fmt.Printf("\nРезультаты вычислений:\n")
	fmt.Printf("Определитель: %.2f\n", result.Determinant)
	fmt.Printf("След: %.2f\n", result.Trace)
	fmt.Println("Транспонированная матрица:")
	for _, row := range result.Transposed {
		fmt.Println(row)
	}
}

func main() {
	matrixOps := MatrixOperations{}

	if err := matrixOps.LoadDataFromFile("input.txt"); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	if err := matrixOps.ValidateSquareMatrix(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	results := matrixOps.ComputeMatrixProperties()

	if err := SaveMatrixResults("output.txt", results); err != nil {
		fmt.Printf("Ошибка сохранения результатов: %v\n", err)
		return
	}

	DisplayResults(results)
	fmt.Println("\nРезультаты успешно сохранены в output.txt")
}
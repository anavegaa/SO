package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Lee la matriz de un archivo sin línea de dimensiones (detecta dimensiones automáticamente)
func readMatrix(filename string) ([][]float64, int, int) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var matrix [][]float64
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) == 0 {
			continue
		}
		row := make([]float64, len(line))
		for i, val := range line {
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				panic(err)
			}
			row[i] = f
		}
		matrix = append(matrix, row)
	}

	if len(matrix) == 0 {
		panic("Matriz vacía")
	}

	rows := len(matrix)
	cols := len(matrix[0])
	return matrix, rows, cols
}

func writeMatrix(filename string, matrix [][]float64) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, row := range matrix {
		for j, val := range row {
			fmt.Fprintf(file, "%.6f", val) // escribe con 6 decimales
			if j < len(row)-1 {
				fmt.Fprint(file, " ")
			}
		}
		fmt.Fprintln(file)
	}
}

func multiplySequential(A [][]float64, B [][]float64) [][]float64 {
	N, M := len(A), len(A[0])
	P := len(B[0])

	C := make([][]float64, N)
	for i := 0; i < N; i++ {
		C[i] = make([]float64, P)
		for j := 0; j < P; j++ {
			sum := 0.0
			for k := 0; k < M; k++ {
				sum += A[i][k] * B[k][j]
			}
			C[i][j] = sum
		}
	}
	return C
}

func main() {
	A, _, _ := readMatrix("A.txt")
	B, _, _ := readMatrix("B.txt")

	start := time.Now()
	C := multiplySequential(A, B)
	elapsed := time.Since(start).Seconds()

	writeMatrix("go_seq.txt", C)
	fmt.Printf("Tiempo secuencial: %.4f segundos\n", elapsed)
}

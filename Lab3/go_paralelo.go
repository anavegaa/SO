package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func readMatrix(filename string) ([][]int, int, int) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	dims := strings.Fields(scanner.Text())
	rows, _ := strconv.Atoi(dims[0])
	cols, _ := strconv.Atoi(dims[1])

	matrix := make([][]int, rows)
	for i := 0; i < rows; i++ {
		scanner.Scan()
		line := strings.Fields(scanner.Text())
		matrix[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			matrix[i][j], _ = strconv.Atoi(line[j])
		}
	}
	return matrix, rows, cols
}

func writeMatrix(filename string, matrix [][]int) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, row := range matrix {
		for j, val := range row {
			fmt.Fprintf(file, "%d", val)
			if j < len(row)-1 {
				fmt.Fprint(file, " ")
			}
		}
		fmt.Fprintln(file)
	}
}

func multiplyParallel(A [][]int, B [][]int, numWorkers int) [][]int {
	N, M := len(A), len(A[0])
	P := len(B[0])

	C := make([][]int, N)
	for i := range C {
		C[i] = make([]int, P)
	}

	var wg sync.WaitGroup
	rowsPerWorker := N / numWorkers

	for w := 0; w < numWorkers; w++ {
		start := w * rowsPerWorker
		end := start + rowsPerWorker
		if w == numWorkers-1 {
			end = N
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				for j := 0; j < P; j++ {
					sum := 0
					for k := 0; k < M; k++ {
						sum += A[i][k] * B[k][j]
					}
					C[i][j] = sum
				}
			}
		}(start, end)
	}

	wg.Wait()
	return C
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run matrix_par.go <num_procesos>")
		return
	}
	numWorkers, _ := strconv.Atoi(os.Args[1])

	A, _, _ := readMatrix("A.txt")
	B, _, _ := readMatrix("B.txt")

	start := time.Now()
	C := multiplyParallel(A, B, numWorkers)
	elapsed := time.Since(start).Seconds()

	writeMatrix("C_par.txt", C)
	fmt.Printf("Tiempo paralelo (%d procesos): %.4f segundos\n", numWorkers, elapsed)
}

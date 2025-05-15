package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Matrix [][]float64

// Leer matriz desde archivo
func readMatrix(filename string) (Matrix, int, int) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var matrix Matrix
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

// Escribir matriz a archivo
func writeMatrix(filename string, matrix Matrix) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, row := range matrix {
		for j, val := range row {
			fmt.Fprintf(file, "%.6f", val)
			if j < len(row)-1 {
				fmt.Fprint(file, " ")
			}
		}
		fmt.Fprintln(file)
	}
}

// Multiplicación secuencial
func multiplySequential(A, B Matrix) Matrix {
	N, M := len(A), len(A[0])
	P := len(B[0])

	C := make(Matrix, N)
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

// Multiplicación paralela con pipes
func multiplyAndPipe(A, B Matrix, writer *io.PipeWriter) {
	defer writer.Close()

	rowsA := len(A)
	colsA := len(A[0])
	colsB := len(B[0])

	C := make(Matrix, rowsA)
	for i := 0; i < rowsA; i++ {
		C[i] = make([]float64, colsB)
		for j := 0; j < colsB; j++ {
			sum := 0.0
			for k := 0; k < colsA; k++ {
				sum += A[i][k] * B[k][j]
			}
			C[i][j] = sum
		}
	}

	encoder := gob.NewEncoder(writer)
	err := encoder.Encode(C)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing matrix to pipe: %v\n", err)
	}
}

func readFromPipe(reader *io.PipeReader, wg *sync.WaitGroup, outputFile string) {
	defer wg.Done()
	defer reader.Close()

	decoder := gob.NewDecoder(reader)
	var C Matrix
	err := decoder.Decode(&C)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding matrix from pipe: %v\n", err)
		return
	}

	writeMatrix(outputFile, C)
	fmt.Printf("Resultado paralelo guardado en %s\n", outputFile)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run matrix.go [secuencial | paralelo]")
		return
	}

	A, _, _ := readMatrix("A.txt")
	B, _, _ := readMatrix("B.txt")

	if len(A) == 0 || len(B) == 0 || len(A[0]) != len(B) {
		fmt.Println("Matrices no multiplicables: columnas de A deben ser iguales a filas de B.")
		return
	}

	mode := os.Args[1]
	switch mode {
	case "secuencial":
		start := time.Now()
		C := multiplySequential(A, B)
		elapsed := time.Since(start).Seconds()
		writeMatrix("go_seq.txt", C)
		fmt.Printf("Resultado secuencial guardado en go_seq.txt\n")
		fmt.Printf("Tiempo secuencial: %.10f segundos\n", elapsed)

	case "paralelo":
		reader, writer := io.Pipe()
		var wg sync.WaitGroup
		wg.Add(1)

		start := time.Now()
		go multiplyAndPipe(A, B, writer)
		go readFromPipe(reader, &wg, "go_par.txt")
		wg.Wait()
		duration := time.Since(start)
		fmt.Printf("Tiempo paralelo: %.4f segundos\n", duration.Seconds())

	default:
		fmt.Println("Modo no válido. Usa 'secuencial' o 'paralelo'")
	}
}

//Juan DIego Calderon Bermeo 1000378849
//Ana Maria Vega Angarita 1004945529

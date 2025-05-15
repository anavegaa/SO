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

// Lee matriz desde archivo sin encabezado
func readMatrix(path string) (Matrix, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var matrix Matrix
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		values := strings.Fields(line)
		row := make([]float64, len(values))
		for i, val := range values {
			row[i], err = strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing float: %v", err)
			}
		}
		matrix = append(matrix, row)
	}
	return matrix, scanner.Err()
}

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

func readFromPipe(reader *io.PipeReader, wg *sync.WaitGroup) {
	defer wg.Done()
	defer reader.Close()

	decoder := gob.NewDecoder(reader)
	var C Matrix
	err := decoder.Decode(&C)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding matrix from pipe: %v\n", err)
		return
	}

	fmt.Println("Resultado:")
	for _, row := range C {
		for _, val := range row {
			fmt.Printf("%.2f ", val)
		}
		fmt.Println()
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run go_parallel.go A.txt B.txt")
		return
	}

	A, err := readMatrix(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error leyendo matriz A: %v\n", err)
		return
	}
	B, err := readMatrix(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error leyendo matriz B: %v\n", err)
		return
	}

	if len(A) == 0 || len(B) == 0 || len(A[0]) != len(B) {
		fmt.Println("Matrices no multiplicables: columnas de A deben ser iguales a filas de B.")
		return
	}

	reader, writer := io.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)

	start := time.Now()

	go multiplyAndPipe(A, B, writer)
	go readFromPipe(reader, &wg)

	wg.Wait()

	duration := time.Since(start)
	fmt.Printf("Tiempo de ejecución: %v\n", duration)
}

#include <iostream>
#include <fstream>
#include <vector>
#include <sstream>
#include <string>

int main() {
    std::ifstream archivoA("A.txt");
    std::ifstream archivoB("B.txt");

    if (!archivoA.is_open()) {
        std::cerr << "No se pudo abrir el archivo A.txt" << std::endl;
        return 1;
    }
    if (!archivoB.is_open()) {
        std::cerr << "No se pudo abrir el archivo B.txt" << std::endl;
        return 1;
    }

    // Leer matriz A
    std::vector<std::vector<double>> A;
    std::string linea;
    int columnasA = 0;

    while (std::getline(archivoA, linea)) {
        std::stringstream ss(linea);
        double valor;
        std::vector<double> fila;

        while (ss >> valor) {
            fila.push_back(valor);
        }

        if (!fila.empty()) {
            A.push_back(fila);
            if (columnasA == 0) {
                columnasA = fila.size();
            } else if (columnasA != fila.size()) {
                std::cerr << "Filas con diferentes cantidades de columnas en A.txt" << std::endl;
                return 1;
            }
        }
    }

    int filasA = A.size();

    // Leer matriz B
    std::vector<std::vector<double>> B;
    int columnasB = 0;

    while (std::getline(archivoB, linea)) {
        std::stringstream ss(linea);
        double valor;
        std::vector<double> fila;

        while (ss >> valor) {
            fila.push_back(valor);
        }

        if (!fila.empty()) {
            B.push_back(fila);
            if (columnasB == 0) {
                columnasB = fila.size();
            } else if (columnasB != fila.size()) {
                std::cerr << "Filas con diferentes cantidades de columnas en B.txt" << std::endl;
                return 1;
            }
        }
    }

    int filasB = B.size();

    std::cout << "Dimensiones A: " << filasA << "x" << columnasA << std::endl;
    std::cout << "Dimensiones B: " << filasB << "x" << columnasB << std::endl;

    // Verificar compatibilidad
    if (columnasA != filasB) {
        std::cerr << "Las dimensiones no son compatibles para la multiplicación." << std::endl;
        return 1;
    }

    // Multiplicación de matrices
    std::vector<std::vector<double>> C(filasA, std::vector<double>(columnasB, 0.0));

    for (int i = 0; i < filasA; i++) {
        for (int j = 0; j < columnasB; j++) {
            for (int k = 0; k < columnasA; k++) {
                C[i][j] += A[i][k] * B[k][j];
            }
        }
    }

    // Guardar matriz C
    std::ofstream archivoC("C_seq.txt");
    if (archivoC.is_open()) {
        for (const auto& fila : C) {
            for (double val : fila) {
                archivoC << val << " ";
            }
            archivoC << std::endl;
        }
        archivoC.close();
        std::cout << "Matriz C guardada en C_seq.txt" << std::endl;
    } else {
        std::cerr << "No se pudo escribir en C_seq.txt" << std::endl;
        return 1;
    }

    return 0;
}

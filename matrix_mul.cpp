#include <iostream>
#include <fstream>
#include <vector>
#include <sstream>
#include <string>
#include <sys/shm.h>
#include <sys/wait.h>
#include <unistd.h>
#include <ctime>

using namespace std;

vector<vector<double>> readMatrix(const string& filename, int& rows, int& cols) {
    ifstream file(filename);
    string line;
    vector<vector<double>> matrix;
    cols = -1;

    while (getline(file, line)) {
        stringstream ss(line);
        double value;
        vector<double> row;
        while (ss >> value) {
            row.push_back(value);
        }
        if (!row.empty()) {
            if (cols == -1)
                cols = (int)row.size();
            else if (cols != (int)row.size()) {
                cerr << "Error: filas con diferente número de columnas en " << filename << endl;
                exit(1);
            }
            matrix.push_back(row);
        }
    }
    rows = (int)matrix.size();
    return matrix;
}

void writeMatrix(const vector<vector<double>>& matrix, const string& filename) {
    ofstream file(filename);
    for (auto& row : matrix) {
        for (auto val : row) {
            file << val << " ";
        }
        file << endl;
    }
}

void writeMatrix(double* C, int N, int P, const string& filename) {
    ofstream file(filename);
    for (int i = 0; i < N; ++i) {
        for (int j = 0; j < P; ++j)
            file << C[i * P + j] << " ";
        file << endl;
    }
}

// Implementación secuencial (tu primer código)
void sequentialMultiply() {
    int N, M, P;
    auto A = readMatrix("A.txt", N, M);
    auto B = readMatrix("B.txt", M, P);

    if (M != (int)B.size()) {
        cerr << "Dimensiones incompatibles para multiplicacion" << endl;
        exit(1);
    }

    vector<vector<double>> C(N, vector<double>(P, 0.0));

    for (int i = 0; i < N; i++) {
        for (int j = 0; j < P; j++) {
            for (int k = 0; k < M; k++) {
                C[i][j] += A[i][k] * B[k][j];
            }
        }
    }

    writeMatrix(C, "C_seq.txt");
    cout << "Matriz C secuencial guardada en C_seq.txt" << endl;
}

// Implementación paralela (tu segundo código)
void parallelMultiply(int K) {
    int N, M, P;
    auto A = readMatrix("A.txt", N, M);
    auto B = readMatrix("B.txt", M, P);

    if (M != (int)B.size()) {
        cerr << "Dimensiones incompatibles para multiplicacion" << endl;
        exit(1);
    }

    // Crear memoria compartida para C
    int shmid = shmget(IPC_PRIVATE, sizeof(double) * N * P, IPC_CREAT | 0666);
    if (shmid < 0) {
        perror("shmget");
        exit(1);
    }
    double* C = (double*)shmat(shmid, NULL, 0);
    if (C == (void*) -1) {
        perror("shmat");
        exit(1);
    }

    clock_t start = clock();

    int rows_per_proc = N / K;
    for (int i = 0; i < K; i++) {
        int start_row = i * rows_per_proc;
        int end_row = (i == K - 1) ? N : start_row + rows_per_proc;

        pid_t pid = fork();
        if (pid < 0) {
            perror("fork");
            exit(1);
        }

        if (pid == 0) {  // Proceso hijo calcula parte de la matriz
            for (int r = start_row; r < end_row; r++) {
                for (int j = 0; j < P; j++) {
                    double sum = 0;
                    for (int k = 0; k < M; k++) {
                        sum += A[r][k] * B[k][j];
                    }
                    C[r * P + j] = sum;
                }
            }
            _exit(0);
        }
    }

    // Esperar a todos los hijos
    for (int i = 0; i < K; i++) {
        wait(NULL);
    }

    clock_t end = clock();

    writeMatrix(C, N, P, "C_par.txt");
    cout << "Matriz C paralela guardada en C_par.txt" << endl;
    cout << "Tiempo paralelo (" << K << " procesos): " << double(end - start) / CLOCKS_PER_SEC << " segundos" << endl;

    shmdt(C);
    shmctl(shmid, IPC_RMID, NULL);
}

int main(int argc, char* argv[]) {
    if (argc < 2) {
        cout << "Uso: " << argv[0] << " [secuencial | paralelo <num_procesos>]" << endl;
        return 1;
    }

    string modo = argv[1];
    if (modo == "secuencial") {
        sequentialMultiply();
    } else if (modo == "paralelo") {
        if (argc != 3) {
            cerr << "Especifica el número de procesos: " << argv[0] << " paralelo <num_procesos>" << endl;
            return 1;
        }
        int num_procesos = stoi(argv[2]);
        parallelMultiply(num_procesos);
    } else {
        cerr << "Modo no reconocido. Usa 'secuencial' o 'paralelo'" << endl;
        return 1;
    }

    return 0;
}
// This code implements matrix multiplication in both sequential and parallel modes.
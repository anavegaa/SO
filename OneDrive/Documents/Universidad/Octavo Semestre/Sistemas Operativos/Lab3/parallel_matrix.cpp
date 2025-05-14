// parallel_matrix.cpp
#include <iostream>
#include <fstream>
#include <vector>
#include <sstream>
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
                cols = row.size();
            else if (cols != (int)row.size()) {
                cerr << "Error: Filas con diferente número de columnas en " << filename << endl;
                exit(1);
            }
            matrix.push_back(row);
        }
    }

    rows = matrix.size();
    return matrix;
}

void writeMatrix(double* C, int N, int P, const string& filename) {
    ofstream file(filename);
    for (int i = 0; i < N; ++i) {
        for (int j = 0; j < P; ++j)
            file << C[i * P + j] << " ";
        file << endl;
    }
}

int main(int argc, char* argv[]) {
    if (argc != 2) {
        cerr << "Usage: ./parallel_matrix <num_processes>\n";
        return 1;
    }

    int K = stoi(argv[1]);
    int N, M, P;

    // Leer matrices A y B sin dimensiones en la primera línea
    auto A = readMatrix("A.txt", N, M);
    auto B = readMatrix("B.txt", M, P);

    // Verificar compatibilidad
    if ((int)B.size() != M) {
        cerr << "Error: Columnas de A no coinciden con filas de B." << endl;
        return 1;
    }

    // Shared memory para C
    int shmid = shmget(IPC_PRIVATE, sizeof(double) * N * P, IPC_CREAT | 0666);
    double* C = (double*)shmat(shmid, NULL, 0);

    clock_t start = clock();

    int rows_per_proc = N / K;
    for (int i = 0; i < K; ++i) {
        int start_row = i * rows_per_proc;
        int end_row = (i == K - 1) ? N : start_row + rows_per_proc;

        if (fork() == 0) {
            for (int r = start_row; r < end_row; ++r)
                for (int j = 0; j < P; ++j) {
                    C[r * P + j] = 0;
                    for (int k = 0; k < M; ++k)
                        C[r * P + j] += A[r][k] * B[k][j];
                }
            _exit(0);
        }
    }

    for (int i = 0; i < K; ++i)
        wait(NULL);
    clock_t end = clock();

    writeMatrix(C, N, P, "C.txt");
    cout << "Parallel time (" << K << " processes): " << double(end - start) / CLOCKS_PER_SEC << " seconds\n";

    shmdt(C);
    shmctl(shmid, IPC_RMID, NULL);
    return 0;
}

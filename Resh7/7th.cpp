#include <iostream>
#include <fstream>
using namespace std;

void freeMatrix(int** matrix, int rows) {
    if (matrix != nullptr) {
        for (int i = 0; i < rows; i++) {
            delete[] matrix[i];
        }
        delete[] matrix;
    }
}

void readMatrix(ifstream& in, int**& matrix, int& rows, int& cols) {
    in >> rows >> cols;
    matrix = new int*[rows];
    for (int i = 0; i < rows; i++) {
        matrix[i] = new int[cols];
        for (int j = 0; j < cols; j++) {
            in >> matrix[i][j];
        }
    }
}

void multiply(int** a, int** b, int**& result, int r1, int c1, int r2, int c2) {
    result = new int*[r1];
    for (int i = 0; i < r1; i++) {
        result[i] = new int[c2];
        for (int j = 0; j < c2; j++) {
            result[i][j] = 0;
            for (int k = 0; k < c1; k++) {
                result[i][j] += a[i][k] * b[k][j];
            }
        }
    }
}

void writeMatrix(ofstream& out, int** matrix, int rows, int cols) {
    for (int i = 0; i < rows; i++) {
        for (int j = 0; j < cols; j++) {
            out << matrix[i][j] << " ";
        }
        out << endl;
    }
}

int main() {
    ifstream in("C:\\c++\\input.txt");
    ofstream out("C:\\c++\\output.txt");

    if (!in.is_open() || !out.is_open()) {
        cout << "n/a" << endl;
        return 1;
    }

    int r1, c1, r2, c2;
    int** mat1 = nullptr, ** mat2 = nullptr, ** res = nullptr;

    readMatrix(in, mat1, r1, c1);
    readMatrix(in, mat2, r2, c2);

    if (c1 != r2) {
        cout << "n/a" << endl;
        freeMatrix(mat1, r1);
        freeMatrix(mat2, r2);
        return 1;
    }

    multiply(mat1, mat2, res, r1, c1, r2, c2);
    writeMatrix(out, res, r1, c2);

    freeMatrix(mat1, r1);
    freeMatrix(mat2, r2);
    freeMatrix(res, r1);

    return 0;
}
#include <iostream>
#include <fstream>
#include <vector>

using namespace std;

vector<vector<double>> readMatrix() {
    ifstream fin("input.txt");
    vector<vector<double>> matrix;

    if (!fin.is_open()) {
        cout << "Не удалось открыть файл!" << endl;
        return matrix;
    }

    double num;
    // Читаем числа пока они есть в файле
    while (fin >> num) {
        vector<double> row;
        row.push_back(num);

        // Читаем остальные числа в строке
        while (fin.peek() != '\n' && fin.peek() != EOF) {
            fin >> num;
            row.push_back(num);
        }

        matrix.push_back(row);
    }

    fin.close();
    return matrix;
}

double calculateTrace(const vector<vector<double>>& mat) {
    double trace = 0;
    for (int i = 0; i < mat.size(); i++) {
        trace += mat[i][i];
    }
    return trace;
}

vector<vector<double>> transposeMatrix(const vector<vector<double>>& mat) {
    vector<vector<double>> result(mat[0].size(), vector<double>(mat.size()));

    for (int i = 0; i < mat.size(); i++) {
        for (int j = 0; j < mat[0].size(); j++) {
            result[j][i] = mat[i][j];
        }
    }

    return result;
}

double calculateDeterminant(vector<vector<double>> mat) {
    int n = mat.size();

    if (n == 1) return mat[0][0];
    if (n == 2) {
        return mat[0][0] * mat[1][1] - mat[0][1] * mat[1][0];
    }

    double det = 0;

    for (int col = 0; col < n; col++) {
        vector<vector<double>> submat(n - 1, vector<double>(n - 1));

        for (int i = 1; i < n; i++) {
            int subcol = 0;
            for (int j = 0; j < n; j++) {
                if (j == col) continue;
                submat[i - 1][subcol] = mat[i][j];
                subcol++;
            }
        }

        double sign = (col % 2 == 0) ? 1 : -1;
        det += sign * mat[0][col] * calculateDeterminant(submat);
    }

    return det;
}

void writeResults(double trace, double det, const vector<vector<double>>& transposed) {
    ofstream fout("output.txt");

    fout << "След матрицы: " << trace << endl;
    fout << "Определитель матрицы: " << det << endl;
    fout << "Транспонированная матрица:" << endl;

    for (const auto& row : transposed) {
        for (double num : row) {
            fout << num << " ";
        }
        fout << endl;
    }

    fout.close();
}

int main() {
    vector<vector<double>> matrix = readMatrix();

    if (matrix.empty() || matrix.size() != matrix[0].size()) {
        cout << "ERROR: matrix not square" << endl;
        return 1;
    }

    double trace = calculateTrace(matrix);
    double det = calculateDeterminant(matrix);
    vector<vector<double>> transposed = transposeMatrix(matrix);

    writeResults(trace, det, transposed);

    cout << "Result in output.txt" << endl;

    return 0;
}
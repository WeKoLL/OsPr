import copy

def read_matrix(filename):
    with open(filename, "r") as file:
        return [[float(x) for x in line.strip().split()] for line in file]

def write_results(filename, det, trace, transposed):
    with open(filename, 'w') as f:
        f.write(f"Determinant: {det}\nTrace: {trace}\nTransposed:\n")
        f.write('\n'.join(' '.join(map(str, row)) for row in transposed))

def get_trace(matrix):
    return sum(matrix[i][i] for i in range(len(matrix)))

def transpose_matrix(matrix):
    return [[matrix[i][j] for j in range(len(matrix[0]))] for i in range(len(matrix))]

def get_submatrix(matrix, exclude_row, exclude_col):
    return [[matrix[i][j] for j in range(len(matrix[i])) if j != exclude_col] 
            for i in range(len(matrix)) if i != exclude_row]

def compute_determinant(matrix):
    n = len(matrix)
    if n == 1:
        return matrix[0][0]
    if n == 2:
        return matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]
    return sum((-1)**j * matrix[0][j] * compute_determinant(get_submatrix(matrix, 0, j)) 
               for j in range(n))

def main():
    matrix = read_matrix("input.txt")
    if any(len(row) != len(matrix) for row in matrix):
        with open("output.txt", "w") as f:
            f.write("Ошибка: матрица не квадратная\n")
        return
    
    det = compute_determinant(matrix)
    trace = get_trace(matrix)
    transposed = transpose_matrix(matrix)
    write_results("output.txt", det, trace, transposed)

if __name__ == "__main__":
    main()
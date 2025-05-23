from collections import Counter

def read_numbers(filename):
    with open(filename, 'r') as f:
        return list(map(int, f.read().split()))

def process_numbers(a, b):
    count_a, count_b = Counter(a), Counter(b)
    common = set(a) & set(b)
    result = []
    for num in common:
        result.extend([num] * max(count_a[num], count_b[num]))
    return sorted(result)

def write_results(filename, numbers):
    with open(filename, 'w') as f:
        f.write('\n'.join(map(str, numbers)))

def main():
    a = read_numbers('A.txt')
    b = read_numbers('B.txt')
    write_results('C.txt', process_numbers(a, b))

if __name__ == "__main__":
    main()
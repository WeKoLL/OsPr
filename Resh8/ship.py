
from string import ascii_uppercase

file = (open("input.txt").readline()).upper()

n = int(input())

def Cez(n, file):
    res = ''
    l =ascii_uppercase
    for i in file:
        res += l[(l.index(i) + n) % len(l)]
    return res
def Atbash(file):
    l = ascii_uppercase
    return file.translate(str.maketrans(l + l.upper(), l[::-1] + l.upper()[::-1]))

print(Atbash(file))
print(Cez(n, file))
package main

import (
	"fmt"
	"os"
)

func atbash(s string) string {
	output := make([]rune, len(s))
	for i, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			output[i] = 'Z' - (r - 'A')
		case r >= 'a' && r <= 'z':
			output[i] = 'z' - (r - 'a')
		default:
			output[i] = r
		}
	}
	return string(output)
}

func caesar(s string, shift int) string {
	if shift < 0 {
		shift = 26 + shift%26
	}
	
	output := make([]rune, len(s))
	for i, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			output[i] = 'A' + (r-'A'+rune(shift))%26
		case r >= 'a' && r <= 'z':
			output[i] = 'a' + (r-'a'+rune(shift))%26
		default:
			output[i] = r
		}
	}
	return string(output)
}

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("Не удалось прочитать файл:", err)
		return
	}
	text := string(content)

	var n int
	fmt.Print("Укажите сдвиг: ")
	if _, err := fmt.Scan(&n); err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}

	a := atbash(text)
	c := caesar(text, n)

	result := fmt.Sprintf("Атбаш:\n%s\n\nЦезарь (сдвиг %d):\n%s", a, n, c)
	if err := os.WriteFile("output.txt", []byte(result), 0644); err != nil {
		fmt.Println("Ошибка записи:", err)
		return
	}

	fmt.Println("Результат сохранен в output.txt")
}

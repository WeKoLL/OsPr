package main

import (
	"fmt"
	"os"
)

func atbashCipher(text string) string {
	encrypted := make([]rune, 0, len(text))
	for _, char := range text {
		switch {
		case 'A' <= char && char <= 'Z':
			encrypted = append(encrypted, 'Z'-(char-'A'))
		case 'a' <= char && char <= 'z':
			encrypted = append(encrypted, 'z'-(char-'a'))
		default:
			encrypted = append(encrypted, char)
		}
	}
	return string(encrypted)
}

func caesarCipher(text string, shift int) string {
	encrypted := make([]rune, 0, len(text))
	shift = (shift%26 + 26) % 26
	for _, char := range text {
		switch {
		case 'A' <= char && char <= 'Z':
			encrypted = append(encrypted, 'A'+(char-'A'+rune(shift))%26)
		case 'a' <= char && char <= 'z':
			encrypted = append(encrypted, 'a'+(char-'a'+rune(shift))%26)
		default:
			encrypted = append(encrypted, char)
		}
	}
	return string(encrypted)
}

func main() {
	inputData, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	var shift int
	fmt.Print("Введите сдвиг для Цезаря: ")
	_, err = fmt.Scan(&shift)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}

	output := fmt.Sprintf(
		"Атбаш: %s\nЦезарь (сдвиг %d): %s",
		atbashCipher(string(inputData)),
		shift,
		caesarCipher(string(inputData), shift),
	)

	err = os.WriteFile("output.txt", []byte(output), 0644)
	if err != nil {
		fmt.Println("Ошибка записи:", err)
		return
	}

	fmt.Println("Результат сохранён в output.txt")
}
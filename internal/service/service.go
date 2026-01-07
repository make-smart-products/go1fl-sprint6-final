package service

import (
	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func Converter(input string) string { // Функция Converter принимает строку и возвращает строку

	for _, char := range input { // Перебор каждого символа во входной строке
		if char != '.' && char != '-' && char != ' ' {
			return morse.ToMorse(input) // Конвертируем всю строку из обычного текста в азбуку Морзе
		}
	}

	return morse.ToText(input) // Если все символы были из множества "-. ", конвертируем из Морзе в текст
}

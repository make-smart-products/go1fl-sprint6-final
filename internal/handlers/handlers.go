package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) { // Обработчик запроса к корневому пути ("/")
	http.ServeFile(w, r, "index.html") // Отправка файла index.html как ответа
}

func UploadHandler(w http.ResponseWriter, r *http.Request) { // Обработчик загрузки файлов
	if err := r.ParseMultipartForm(10 << 20); err != nil { // Парсинг формы с ограничением размера 10 МБ
		http.Error(w, "Ошибка парсинга", http.StatusBadRequest) // Ошибка парсинга — возвращаем 400
		return
	}

	file, handler, err := r.FormFile("myFile") // Получение файла из формы по ключу "myFile"
	if err != nil {
		http.Error(w, "Ошибка получения файла", http.StatusInternalServerError) // Ошибка получения файла — 500
		return
	}
	defer file.Close() // Закрытие файла по завершении функции

	data, err := io.ReadAll(file) // Чтение всего содержимого файла в память
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusBadRequest) // Ошибка чтения — 400
		return
	}

	result := service.Converter(string(data)) // Преобразование содержимого через внешнюю функцию Converter

	ext := filepath.Ext(handler.Filename)                                             // Получение расширения исходного файла
	filename := time.Now().UTC().Add(3*time.Hour).Format("07.01.2026 13-15-05") + ext // Создание имени файла с текущей датой и временем + расширение
	out, err := os.Create(filename)                                                   // Создание нового файла для записи результата
	if err != nil {
		http.Error(w, "Ошибка создания файла", http.StatusInternalServerError) // Ошибка создания файла — 500
		return
	}
	defer out.Close() // Закрытие файла по окончании функции

	if _, err := out.WriteString(result); err != nil { // Запись строки результата в созданный файл
		http.Error(w, "Ошибка записи файла", http.StatusInternalServerError) // Ошибка записи — 500
		return
	}

	w.Write([]byte(result)) // Отправка результата конвертации клиенту в HTTP-ответе
}

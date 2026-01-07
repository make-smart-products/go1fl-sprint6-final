package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

// Server структура сервера с логгером и http.Server
type Server struct {
	Logger *log.Logger
	HTTP   *http.Server
}

func NewServer(logger *log.Logger) *Server { // Функция-конструктор для создания нового сервера
	mux := http.NewServeMux()                         // Создание мультиплексора для маршрутизации HTTP-запросов
	mux.HandleFunc("/", handlers.IndexHandler)        // Регистрация обработчика для корневого пути "/"
	mux.HandleFunc("/upload", handlers.UploadHandler) // Регистрация обработчика для пути "/upload"

	srv := &http.Server{ // Создание экземпляра HTTP-сервера с настройками:
		Addr:         ":8080",          // - Порт прослушивания 8080
		Handler:      mux,              // - Установка маршрутизатора для обработки запросов
		ErrorLog:     logger,           // - Использование переданного логгера для ошибок HTTP-сервера
		ReadTimeout:  5 * time.Second,  // - Тайм-аут на чтение запроса (5 секунд)
		WriteTimeout: 10 * time.Second, // - Тайм-аут на запись ответа (10 секунд)
		IdleTimeout:  15 * time.Second, // - Тайм-аут бездействия соединения (15 секунд)
	}

	return &Server{ // Возврат указателя на созданную структуру Server
		Logger: logger,
		HTTP:   srv,
	}
}

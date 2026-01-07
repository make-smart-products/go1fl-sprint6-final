package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "server: ", log.LstdFlags|log.Lshortfile) // Создание нового логгера, который будет выводить сообщения в консоль, с префиксом "server: " и флагами вывода даты, времени и файла
	serv := server.NewServer(logger)                                       // Создание нового экземпляра сервера с логгером
	logger.Println("Сервер работает на порту :8080")                       // Запись в лог информации о запуске сервера и порте, на котором он работает
	err := serv.HTTP.ListenAndServe()                                      // Запуск HTTP-сервера, метод ListenAndServe блокирует выполнение, пока сервер работает или не возникнет ошибка
	if err != nil {                                                        // Если при запуске сервера возникла ошибка, она логируется и программа завершается
		logger.Fatal(err)
	}
}

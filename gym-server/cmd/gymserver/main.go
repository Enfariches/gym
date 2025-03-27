package main

import (
	"fmt"
	"health/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// Логгер
	// Приложение
	// grpc-сервер
	// Тесты
}

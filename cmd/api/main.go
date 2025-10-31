package main

import (
	"meawle/cmd/api/config"
	"meawle/cmd/api/di"
	"meawle/cmd/api/routes"
	"meawle/cmd/api/server"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()
	logger := config.SetupLogger()

	// Инициализация зависимостей
	deps, err := di.InitializeDependencies(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize dependencies:", err)
	}
	defer deps.DB.Close()

	// Настройка маршрутов
	router := routes.SetupRoutes(
		deps.UserHandler,
		deps.CatBreedHandler,
		deps.AuthMiddleware,
		deps.CatBreedMiddleware,
	)

	// Создание и запуск сервера
	srv := server.NewServer(cfg, router, logger)
	srv.Start()

	// Ожидание graceful shutdown
	server.WaitForShutdown()
	srv.Shutdown()
}

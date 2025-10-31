package di

import (
	"log"

	"meawle/internal/config"
	"meawle/internal/database"
	"meawle/internal/handlers"
	"meawle/internal/middleware"
	"meawle/internal/repositories"
	"meawle/internal/services"
)

// Dependencies содержит все зависимости приложения
type Dependencies struct {
	Config          *config.Config
	Logger          *log.Logger
	DB              *database.Database
	UserRepo        repositories.UserRepository
	CatBreedRepo    repositories.CatBreedRepository
	UserService     *services.UserService
	CatBreedService *services.CatBreedService
	UserHandler     *handlers.UserHandler
	CatBreedHandler *handlers.CatBreedHandler
	AuthMiddleware  *middleware.AuthMiddleware
}

// InitializeDependencies инициализирует все зависимости приложения
func InitializeDependencies(cfg *config.Config, logger *log.Logger) (*Dependencies, error) {
	// Инициализация базы данных
	db, err := database.New(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	// Запуск миграций
	if err := db.RunMigrations("migrations"); err != nil {
		return nil, err
	}

	// Инициализация репозиториев
	userRepo := repositories.NewUserRepository(db)
	catBreedRepo := repositories.NewCatBreedRepository(db)

	// Инициализация сервисов
	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	catBreedService := services.NewCatBreedService(catBreedRepo)

	// Инициализация хэндлеров
	userHandler := handlers.NewUserHandler(userService)
	catBreedHandler := handlers.NewCatBreedHandler(catBreedService)

	// Инициализация middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)

	return &Dependencies{
		Config:          cfg,
		Logger:          logger,
		DB:              db,
		UserRepo:        userRepo,
		CatBreedRepo:    catBreedRepo,
		UserService:     userService,
		CatBreedService: catBreedService,
		UserHandler:     userHandler,
		CatBreedHandler: catBreedHandler,
		AuthMiddleware:  authMiddleware,
	}, nil
}

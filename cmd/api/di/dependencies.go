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
	CatRepo         repositories.CatRepository
	UserService     *services.UserService
	CatBreedService *services.CatBreedService
	CatService      *services.CatService
	UserHandler     *handlers.UserHandler
	CatBreedHandler *handlers.CatBreedHandler
	CatHandler      *handlers.CatHandler
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
	catRepo := repositories.NewCatRepository(db)

	// Инициализация сервисов
	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	catBreedService := services.NewCatBreedService(catBreedRepo)
	catService := services.NewCatService(catRepo)

	// Инициализация хэндлеров
	userHandler := handlers.NewUserHandler(userService)
	catBreedHandler := handlers.NewCatBreedHandler(catBreedService)
	catHandler := handlers.NewCatHandler(catService)

	// Инициализация middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)

	return &Dependencies{
		Config:          cfg,
		Logger:          logger,
		DB:              db,
		UserRepo:        userRepo,
		CatBreedRepo:    catBreedRepo,
		CatRepo:         catRepo,
		UserService:     userService,
		CatBreedService: catBreedService,
		CatService:      catService,
		UserHandler:     userHandler,
		CatBreedHandler: catBreedHandler,
		CatHandler:      catHandler,
		AuthMiddleware:  authMiddleware,
	}, nil
}

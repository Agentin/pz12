// Package main Notes API server.
//
// @title Notes API
// @version 1.0
// @description Учебный REST API для заметок (CRUD).
// @contact.name Backend Course
// @contact.email example@university.ru
// @BasePath /api/v1
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "example.com/notes-api/docs" // Импорт сгенерированной документации
	httpSwagger "github.com/swaggo/http-swagger"
	"example.com/notes-api/internal/http/handlers"
	"example.com/notes-api/internal/repo"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <token>

func main() {
	// Инициализация репозитория
	noteRepo := repo.NewNoteRepoMem()
	
	// Инициализация хэндлеров
	h := &handlers.Handler{Repo: noteRepo}
	
	// Создание роутера
	r := chi.NewRouter()
	
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	// Настройка маршрутов
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/notes", h.ListNotes)
		r.Post("/notes", h.CreateNote)
		r.Get("/notes/{id}", h.GetNote)
		r.Patch("/notes/{id}", h.UpdateNote)
		r.Delete("/notes/{id}", h.DeleteNote)
	})
	
	// Swagger UI
	r.Get("/docs/*", httpSwagger.WrapHandler)
	
	// Запуск сервера
	log.Println("Server started at :8080")
	log.Println("Swagger UI available at http://localhost:8080/docs/index.html")
	log.Fatal(http.ListenAndServe(":8080", r))
}

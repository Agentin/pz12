package core

import "time"

// Note представляет сущность заметки
type Note struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// NoteCreate DTO для создания заметки
type NoteCreate struct {
	Title   string `json:"title" example:"Новая заметка"`
	Content string `json:"content" example:"Текст заметки"`
}

// NoteUpdate DTO для обновления заметки
type NoteUpdate struct {
	Title   *string `json:"title,omitempty" example:"Обновлённый заголовок"`
	Content *string `json:"content,omitempty" example:"Обновлённый текст"`
}

// PaginatedNotes для пагинированного ответа
type PaginatedNotes struct {
	Notes      []Note `json:"notes"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"totalPages"`
}

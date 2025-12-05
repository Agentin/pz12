package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "github.com/go-chi/chi/v5"
    "example.com/notes-api/internal/core"
    "example.com/notes-api/internal/repo"
)

type Handler struct {
    Repo *repo.NoteRepoMem
}

// ListNotes godoc
// @Summary Список заметок
// @Description Возвращает список заметок с пагинацией и фильтром по заголовку
// @Tags notes
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Размер страницы" default(10) minimum(1) maximum(100)
// @Param q query string false "Поиск по title"
// @Success 200 {array} core.Note
// @Header 200 {integer} X-Total-Count "Общее количество"
// @Failure 500 {object} map[string]string
// @Router /notes [get]
func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
    page := 1
    limit := 10
    query := ""
    
    if pageStr := r.URL.Query().Get("page"); pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }
    
    if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
            limit = l
        }
    }
    
    if q := r.URL.Query().Get("q"); q != "" {
        query = strings.TrimSpace(q)
    }
    
    notes, total, err := h.Repo.GetAll(page, limit, query)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Total-Count", strconv.Itoa(total))
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(notes)
}

// CreateNote godoc
// @Summary Создать заметку
// @Tags notes
// @Accept json
// @Produce json
// @Param input body core.NoteCreate true "Данные новой заметки"
// @Success 201 {object} core.Note
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notes [post]
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
    var input core.NoteCreate
    
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
        return
    }
    
    if input.Title == "" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Title is required"})
        return
    }
    
    note := core.Note{
        Title:   input.Title,
        Content: input.Content,
    }
    
    id, err := h.Repo.Create(note)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
    
    // Получаем созданную заметку
    createdNote, err := h.Repo.GetByID(id)
    if err != nil || createdNote == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve created note"})
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdNote)
}

// GetNote godoc
// @Summary Получить заметку
// @Tags notes
// @Param id path int true "ID"
// @Success 200 {object} core.Note
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notes/{id} [get]
func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
        return
    }
    
    note, err := h.Repo.GetByID(id)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
    
    if note == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"error": "Note not found"})
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(note)
}

// UpdateNote godoc
// @Summary Обновить заметку (частично)
// @Tags notes
// @Accept json
// @Param id path int true "ID"
// @Param input body core.NoteUpdate true "Поля для обновления"
// @Success 200 {object} core.Note
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notes/{id} [patch]
func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
        return
    }
    
    var input core.NoteUpdate
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
        return
    }
    
    // Получаем существующую заметку
    existingNote, err := h.Repo.GetByID(id)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
    
    if existingNote == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"error": "Note not found"})
        return
    }
    
    // Обновляем только переданные поля
    if input.Title != nil {
        existingNote.Title = *input.Title
    }
    if input.Content != nil {
        existingNote.Content = *input.Content
    }
    
    // Сохраняем обновлённую заметку
    err = h.Repo.Update(id, *existingNote)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
    
    // Получаем обновлённую заметку
    updatedNote, err := h.Repo.GetByID(id)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(updatedNote)
}

// DeleteNote godoc
// @Summary Удалить заметку
// @Tags notes
// @Param id path int true "ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /notes/{id} [delete]
func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
        return
    }
    
    // Проверяем существование заметки
    note, err := h.Repo.GetByID(id)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
    
    if note == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"error": "Note not found"})
        return
    }
    
    // Удаляем заметку
    err = h.Repo.Delete(id)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

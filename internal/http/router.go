package httpx

import (
    "github.com/go-chi/chi/v5"
    "example.com/notes-api/internal/http/handlers"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
    r := chi.NewRouter()
    
    r.Route("/api/v1/notes", func(r chi.Router) {
        r.Get("/", h.ListNotes)          // GET /api/v1/notes
        r.Post("/", h.CreateNote)        // POST /api/v1/notes
        r.Get("/{id}", h.GetNote)        // GET /api/v1/notes/{id}
        r.Patch("/{id}", h.PatchNote)    // PATCH /api/v1/notes/{id}
        r.Delete("/{id}", h.DeleteNote)  // DELETE /api/v1/notes/{id}
    })
    
    return r
}

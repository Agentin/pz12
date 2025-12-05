package repo

import (
    "sync"
    "time"
    "example.com/notes-api/internal/core"
)

type NoteRepoMem struct {
    mu    sync.RWMutex
    notes map[int64]*core.Note
    next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
    return &NoteRepoMem{
        notes: make(map[int64]*core.Note),
        next:  1,
    }
}

func (r *NoteRepoMem) Create(n core.Note) (int64, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    n.ID = r.next
    n.CreatedAt = time.Now()
    r.notes[n.ID] = &n
    r.next++
    
    return n.ID, nil
}

func (r *NoteRepoMem) GetAll(page, limit int, query string) ([]core.Note, int, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var filtered []*core.Note
    for _, note := range r.notes {
        if query == "" {
            filtered = append(filtered, note)
        } else if contains(note.Title, query) {
            filtered = append(filtered, note)
        }
    }
    
    total := len(filtered)
    
    // Пагинация
    start := (page - 1) * limit
    if start > total {
        return []core.Note{}, total, nil
    }
    
    end := start + limit
    if end > total {
        end = total
    }
    
    result := make([]core.Note, 0, end-start)
    for i := start; i < end; i++ {
        result = append(result, *filtered[i])
    }
    
    return result, total, nil
}

func (r *NoteRepoMem) GetByID(id int64) (*core.Note, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    note, exists := r.notes[id]
    if !exists {
        return nil, nil
    }
    
    return note, nil
}

func (r *NoteRepoMem) Update(id int64, n core.Note) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    existing, exists := r.notes[id]
    if !exists {
        return nil
    }
    
    n.ID = id
    n.CreatedAt = existing.CreatedAt
    now := time.Now()
    n.UpdatedAt = &now
    
    r.notes[id] = &n
    return nil
}

func (r *NoteRepoMem) Delete(id int64) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.notes[id]; !exists {
        return nil
    }
    
    delete(r.notes, id)
    return nil
}

func contains(str, substr string) bool {
    if substr == "" {
        return true
    }
    for i := 0; i <= len(str)-len(substr); i++ {
        if str[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}

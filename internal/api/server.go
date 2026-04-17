package api

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/gicket/gicket/internal/model"
	"github.com/gicket/gicket/internal/store"
	"github.com/gicket/gicket/web"
)

type Server struct {
	store *store.Store
	mux   *http.ServeMux
}

func NewServer(s *store.Store) *Server {
	srv := &Server{store: s}
	srv.mux = http.NewServeMux()
	srv.routes()
	return srv
}

func (s *Server) Handler() http.Handler {
	return corsMiddleware(s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/tickets", s.handleTickets)
	s.mux.HandleFunc("/api/tickets/", s.handleTicket)

	// Serve embedded WebUI
	webFS, _ := fs.Sub(web.Assets, ".")
	s.mux.Handle("/", http.FileServer(http.FS(webFS)))
}

// GET /api/tickets?status=open|closed|in-progress|all
// POST /api/tickets
func (s *Server) handleTickets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listTickets(w, r)
	case http.MethodPost:
		s.createTicket(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /api/tickets/{id}
// PUT /api/tickets/{id}
// DELETE /api/tickets/{id}
// POST /api/tickets/{id}/comments
func (s *Server) handleTicket(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/tickets/")
	parts := strings.SplitN(path, "/", 2)
	id := parts[0]

	if len(parts) == 2 && parts[1] == "comments" {
		if r.Method == http.MethodPost {
			s.addComment(w, r, id)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getTicket(w, r, id)
	case http.MethodPut:
		s.updateTicket(w, r, id)
	case http.MethodDelete:
		s.deleteTicket(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listTickets(w http.ResponseWriter, r *http.Request) {
	statusParam := r.URL.Query().Get("status")
	var filter model.Status
	if statusParam != "" && statusParam != "all" {
		filter = model.Status(statusParam)
	}

	tickets, err := s.store.List(filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tickets == nil {
		tickets = []*model.Ticket{}
	}
	writeJSON(w, http.StatusOK, tickets)
}

func (s *Server) getTicket(w http.ResponseWriter, _ *http.Request, id string) {
	ticket, err := s.store.Load(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ticket)
}

type createTicketRequest struct {
	Title       string   `json:"title"`
	Priority    string   `json:"priority"`
	Assignee    string   `json:"assignee"`
	Labels      []string `json:"labels"`
	Description string   `json:"description"`
}

func (s *Server) createTicket(w http.ResponseWriter, r *http.Request) {
	var req createTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	now := time.Now()
	ticket := &model.Ticket{
		ID:          store.GenerateID(),
		Title:       req.Title,
		Status:      model.StatusOpen,
		Priority:    model.Priority(req.Priority),
		Assignee:    req.Assignee,
		Labels:      req.Labels,
		Description: req.Description,
		Created:     now,
		Updated:     now,
		Author:      "web",
	}
	if ticket.Priority == "" {
		ticket.Priority = model.PriorityMedium
	}

	if err := s.store.Save(ticket); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, ticket)
}

type updateTicketRequest struct {
	Title       *string  `json:"title"`
	Status      *string  `json:"status"`
	Priority    *string  `json:"priority"`
	Assignee    *string  `json:"assignee"`
	Labels      []string `json:"labels"`
	Description *string  `json:"description"`
}

func (s *Server) updateTicket(w http.ResponseWriter, r *http.Request, id string) {
	ticket, err := s.store.Load(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	var req updateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title != nil {
		ticket.Title = *req.Title
	}
	if req.Status != nil {
		ticket.Status = model.Status(*req.Status)
	}
	if req.Priority != nil {
		ticket.Priority = model.Priority(*req.Priority)
	}
	if req.Assignee != nil {
		ticket.Assignee = *req.Assignee
	}
	if req.Labels != nil {
		ticket.Labels = req.Labels
	}
	if req.Description != nil {
		ticket.Description = *req.Description
	}

	if err := s.store.Save(ticket); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ticket)
}

func (s *Server) deleteTicket(w http.ResponseWriter, _ *http.Request, id string) {
	if err := s.store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Deleted"})
}

type addCommentRequest struct {
	Body string `json:"body"`
}

func (s *Server) addComment(w http.ResponseWriter, r *http.Request, id string) {
	ticket, err := s.store.Load(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	var req addCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Body == "" {
		writeError(w, http.StatusBadRequest, "Comment body is required")
		return
	}

	comment := model.Comment{
		Author: "web",
		Date:   time.Now(),
		Body:   req.Body,
	}
	ticket.Comments = append(ticket.Comments, comment)

	if err := s.store.Save(ticket); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, ticket)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func StartServer(repoPath string, port int) error {
	s, err := store.NewStore(repoPath)
	if err != nil {
		return err
	}

	srv := NewServer(s)
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("gicket server starting on http://localhost%s\n", addr)
	fmt.Printf("Repository: %s\n", repoPath)
	return http.ListenAndServe(addr, srv.Handler())
}

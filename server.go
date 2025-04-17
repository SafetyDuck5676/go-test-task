package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Server представляет HTTP-сервер.
type Server struct {
	queueManager *QueueManager
	port         string
	timeout      time.Duration
}

// NewServer создает новый HTTP-сервер.
func NewServer(queueManager *QueueManager, port string, timeout time.Duration) *Server {
	return &Server{
		queueManager: queueManager,
		port:         port,
		timeout:      timeout,
	}
}

// Start запускает сервер.
func (s *Server) Start() error {
	http.HandleFunc("/queue/", s.handleQueue)
	return http.ListenAndServe(":"+s.port, nil)
}

// handleQueue обрабатывает запросы к очередям.
func (s *Server) handleQueue(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	queueName := parts[2]

	switch r.Method {
	case http.MethodPut:
		s.handlePut(w, r, queueName)
	case http.MethodGet:
		s.handleGet(w, r, queueName)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// handlePut обрабатывает PUT-запросы.
func (s *Server) handlePut(w http.ResponseWriter, r *http.Request, queueName string) {
	var body struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Message == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	queue, err := s.queueManager.GetOrCreateQueue(queueName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	if err := queue.Add(body.Message); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handleGet обрабатывает GET-запросы.
func (s *Server) handleGet(w http.ResponseWriter, r *http.Request, queueName string) {
	queue, err := s.queueManager.GetOrCreateQueue(queueName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	timeout := s.timeout
	if t := r.URL.Query().Get("timeout"); t != "" {
		if parsedTimeout, err := strconv.Atoi(t); err == nil {
			timeout = time.Duration(parsedTimeout) * time.Second
		}
	}

	message, err := queue.Get(timeout)
	if err != nil {
		if err.Error() == "timeout waiting for message" {
			http.Error(w, "timeout", http.StatusNotFound)
		} else {
			http.Error(w, "no messages available", http.StatusNotFound)
		}
		return
	}

	fmt.Fprint(w, message)
}

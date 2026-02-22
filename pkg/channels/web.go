package channels

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/database"
	"github.com/sipeed/picoclaw/pkg/logger"
)

// WebChannel provides a browser-based chat UI for PicoClaw.
type WebChannel struct {
	*BaseChannel
	mu      sync.RWMutex
	clients map[string]chan string // chatID -> SSE channel
	mux     *http.ServeMux
	db      *database.DB // nil if no database configured
	model   string       // configured model name for status reporting
}

func NewWebChannel(messageBus *bus.MessageBus, mux *http.ServeMux, db *database.DB, model string) (*WebChannel, error) {
	wc := &WebChannel{
		BaseChannel: NewBaseChannel("web", nil, messageBus, nil),
		clients:     make(map[string]chan string),
		mux:         mux,
		db:          db,
		model:       model,
	}

	wc.registerRoutes()
	return wc, nil
}

func (wc *WebChannel) registerRoutes() {
	// Serve frontend
	wc.mux.HandleFunc("/", wc.handleIndex)
	// Chat API
	wc.mux.HandleFunc("/api/chat", wc.handleChat)
	// SSE stream for responses
	wc.mux.HandleFunc("/api/stream", wc.handleStream)
	// Conversations REST API
	wc.mux.HandleFunc("/api/conversations", wc.handleConversations)
	wc.mux.HandleFunc("/api/conversations/", wc.handleConversationByID)
	// Agent status
	wc.mux.HandleFunc("/api/status", wc.handleStatus)
}

func (wc *WebChannel) Start(ctx context.Context) error {
	wc.setRunning(true)
	logger.InfoC("web", "Web channel started")
	return nil
}

func (wc *WebChannel) Stop(ctx context.Context) error {
	wc.setRunning(false)
	wc.mu.Lock()
	for id, ch := range wc.clients {
		close(ch)
		delete(wc.clients, id)
	}
	wc.mu.Unlock()
	logger.InfoC("web", "Web channel stopped")
	return nil
}

func (wc *WebChannel) Send(ctx context.Context, msg bus.OutboundMessage) error {
	// Persist bot response to database
	if wc.db != nil {
		wc.db.SaveMessage(ctx, msg.ChatID, "assistant", msg.Content)
	}

	wc.mu.RLock()
	ch, ok := wc.clients[msg.ChatID]
	wc.mu.RUnlock()

	if !ok {
		logger.WarnCF("web", "No SSE client for chat", map[string]any{"chat_id": msg.ChatID})
		return nil
	}

	select {
	case ch <- msg.Content:
	case <-time.After(5 * time.Second):
		logger.WarnCF("web", "SSE send timeout", map[string]any{"chat_id": msg.ChatID})
	}
	return nil
}

func (wc *WebChannel) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(webUIHTML))
}

func (wc *WebChannel) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Message string `json:"message"`
		ChatID  string `json:"chat_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Bad request: %v", err), http.StatusBadRequest)
		return
	}
	if req.Message == "" {
		http.Error(w, "Empty message", http.StatusBadRequest)
		return
	}
	if req.ChatID == "" {
		req.ChatID = fmt.Sprintf("web-%d", time.Now().UnixNano())
	}

	// Persist to database
	if wc.db != nil {
		ctx := r.Context()
		// Auto-create conversation with first message as title
		title := req.Message
		if len(title) > 60 {
			title = title[:60] + "..."
		}
		wc.db.CreateConversation(ctx, req.ChatID, title, "web")
		wc.db.SaveMessage(ctx, req.ChatID, "user", req.Message)
	}

	wc.HandleMessage("web-user", req.ChatID, req.Message, nil, nil)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"chat_id": req.ChatID,
	})
}

func (wc *WebChannel) handleStream(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		http.Error(w, "Missing chat_id", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ch := make(chan string, 10)
	wc.mu.Lock()
	wc.clients[chatID] = ch
	wc.mu.Unlock()

	defer func() {
		wc.mu.Lock()
		delete(wc.clients, chatID)
		wc.mu.Unlock()
	}()

	// Send initial connected event
	fmt.Fprintf(w, "data: {\"type\":\"connected\"}\n\n")
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(map[string]string{
				"type":    "message",
				"content": msg,
			})
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

// handleConversations handles GET /api/conversations (list) and POST (create).
func (wc *WebChannel) handleConversations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if wc.db == nil {
		json.NewEncoder(w).Encode(map[string]any{"conversations": []any{}, "db": false})
		return
	}

	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		convos, err := wc.db.ListConversations(ctx, 50, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if convos == nil {
			convos = []database.Conversation{}
		}
		json.NewEncoder(w).Encode(map[string]any{"conversations": convos, "db": true})

	case http.MethodPost:
		var req struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		if req.ID == "" {
			req.ID = fmt.Sprintf("web-%d", time.Now().UnixNano())
		}
		if req.Title == "" {
			req.Title = "New Chat"
		}
		convo, err := wc.db.CreateConversation(ctx, req.ID, req.Title, "web")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(convo)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleConversationByID handles /api/conversations/{id} and /api/conversations/{id}/messages.
func (wc *WebChannel) handleConversationByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if wc.db == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	// Parse path: /api/conversations/{id} or /api/conversations/{id}/messages
	path := strings.TrimPrefix(r.URL.Path, "/api/conversations/")
	parts := strings.SplitN(path, "/", 2)
	id := parts[0]
	subResource := ""
	if len(parts) > 1 {
		subResource = parts[1]
	}

	ctx := r.Context()

	if subResource == "messages" {
		// GET messages for conversation
		msgs, err := wc.db.GetMessages(ctx, id, 200)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if msgs == nil {
			msgs = []database.Message{}
		}
		json.NewEncoder(w).Encode(map[string]any{"messages": msgs})
		return
	}

	switch r.Method {
	case http.MethodGet:
		convo, err := wc.db.GetConversation(ctx, id)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(convo)

	case http.MethodPatch:
		var req struct {
			Title string `json:"title"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		if req.Title != "" {
			wc.db.UpdateConversationTitle(ctx, id, req.Title)
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	case http.MethodDelete:
		wc.db.DeleteConversation(ctx, id)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleStatus returns agent status info.
func (wc *WebChannel) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "online",
		"channel": "web",
		"db":      wc.db != nil,
		"model":   wc.model,
	})
}

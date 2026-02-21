package channels

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/logger"
)

// WebChannel provides a browser-based chat UI for PicoClaw.
type WebChannel struct {
	*BaseChannel
	mu       sync.RWMutex
	clients  map[string]chan string // chatID -> SSE channel
	mux      *http.ServeMux
}

func NewWebChannel(messageBus *bus.MessageBus, mux *http.ServeMux) (*WebChannel, error) {
	wc := &WebChannel{
		BaseChannel: NewBaseChannel("web", nil, messageBus, nil),
		clients:     make(map[string]chan string),
		mux:         mux,
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
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if req.Message == "" {
		http.Error(w, "Empty message", http.StatusBadRequest)
		return
	}
	if req.ChatID == "" {
		req.ChatID = fmt.Sprintf("web-%d", time.Now().UnixNano())
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

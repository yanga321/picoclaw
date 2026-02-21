package database

import (
	"context"
	"time"
)

// Conversation represents a chat conversation.
type Conversation struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Channel   string    `json:"channel"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Preview   string    `json:"preview,omitempty"`
}

// Message represents a single message in a conversation.
type Message struct {
	ID             int       `json:"id"`
	ConversationID string    `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

// CreateConversation inserts a new conversation.
func (db *DB) CreateConversation(ctx context.Context, id, title, channel string) (*Conversation, error) {
	now := time.Now()
	_, err := db.conn.ExecContext(ctx,
		`INSERT INTO conversations (id, title, channel, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $4)
		 ON CONFLICT (id) DO NOTHING`,
		id, title, channel, now,
	)
	if err != nil {
		return nil, err
	}
	return &Conversation{ID: id, Title: title, Channel: channel, CreatedAt: now, UpdatedAt: now}, nil
}

// ListConversations returns recent conversations ordered by update time.
func (db *DB) ListConversations(ctx context.Context, limit, offset int) ([]Conversation, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := db.conn.QueryContext(ctx,
		`SELECT c.id, c.title, c.channel, c.created_at, c.updated_at,
		        COALESCE((SELECT content FROM messages WHERE conversation_id = c.id ORDER BY created_at DESC LIMIT 1), '')
		 FROM conversations c
		 ORDER BY c.updated_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convos []Conversation
	for rows.Next() {
		var c Conversation
		if err := rows.Scan(&c.ID, &c.Title, &c.Channel, &c.CreatedAt, &c.UpdatedAt, &c.Preview); err != nil {
			return nil, err
		}
		convos = append(convos, c)
	}
	return convos, rows.Err()
}

// GetConversation returns a single conversation.
func (db *DB) GetConversation(ctx context.Context, id string) (*Conversation, error) {
	var c Conversation
	err := db.conn.QueryRowContext(ctx,
		`SELECT id, title, channel, created_at, updated_at FROM conversations WHERE id = $1`, id,
	).Scan(&c.ID, &c.Title, &c.Channel, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// UpdateConversationTitle updates the title of a conversation.
func (db *DB) UpdateConversationTitle(ctx context.Context, id, title string) error {
	_, err := db.conn.ExecContext(ctx,
		`UPDATE conversations SET title = $2, updated_at = NOW() WHERE id = $1`, id, title)
	return err
}

// DeleteConversation removes a conversation and all its messages.
func (db *DB) DeleteConversation(ctx context.Context, id string) error {
	_, err := db.conn.ExecContext(ctx, `DELETE FROM conversations WHERE id = $1`, id)
	return err
}

// SaveMessage stores a message and touches the conversation update time.
func (db *DB) SaveMessage(ctx context.Context, conversationID, role, content string) (*Message, error) {
	now := time.Now()
	var id int
	err := db.conn.QueryRowContext(ctx,
		`INSERT INTO messages (conversation_id, role, content, created_at)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		conversationID, role, content, now,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	// Touch conversation
	db.conn.ExecContext(ctx,
		`UPDATE conversations SET updated_at = $2 WHERE id = $1`,
		conversationID, now)

	return &Message{ID: id, ConversationID: conversationID, Role: role, Content: content, CreatedAt: now}, nil
}

// GetMessages returns all messages for a conversation.
func (db *DB) GetMessages(ctx context.Context, conversationID string, limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 200
	}
	rows, err := db.conn.QueryContext(ctx,
		`SELECT id, conversation_id, role, content, created_at
		 FROM messages WHERE conversation_id = $1
		 ORDER BY created_at ASC LIMIT $2`,
		conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

// SetAgentState stores a key-value pair for agent state.
func (db *DB) SetAgentState(ctx context.Context, key, value string) error {
	_, err := db.conn.ExecContext(ctx,
		`INSERT INTO agent_state (key, value, updated_at)
		 VALUES ($1, $2, NOW())
		 ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = NOW()`,
		key, value)
	return err
}

// GetAgentState retrieves a state value by key.
func (db *DB) GetAgentState(ctx context.Context, key string) (string, error) {
	var value string
	err := db.conn.QueryRowContext(ctx,
		`SELECT value FROM agent_state WHERE key = $1`, key).Scan(&value)
	return value, err
}

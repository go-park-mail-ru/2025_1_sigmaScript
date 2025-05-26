package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	kinolkHostEnv = "KINOLK_FRONTEND_HOST"
)

type NotificationData struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
	Date  string `json:"date,omitempty"`
}

type Notification struct {
	Type string           `json:"type"`
	Data NotificationData `json:"data,omitempty"`
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type NotificationServiceInterface interface {
}

type NotificationHandler struct {
	upgrader      websocket.Upgrader
	searchService NotificationServiceInterface
	clients       map[*Client]bool
	clientsMu     sync.Mutex
	sysLogger     *zerolog.Logger
	stop          chan struct{}
}

func NewNotificationHandler(ctx context.Context, notificationService NotificationServiceInterface) *NotificationHandler {
	newNotificationHandler := &NotificationHandler{
		searchService: notificationService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return r.Header.Get("Origin") == viper.GetString(kinolkHostEnv) || true
			},
		},
		clients:   make(map[*Client]bool),
		sysLogger: log.Ctx(ctx),
		stop:      make(chan struct{}),
	}

	go newNotificationHandler.broadcastLoop(newNotificationHandler.stop)

	return newNotificationHandler
}

func (h *NotificationHandler) WSHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msgf("New WS Connection: %v", r.Host)

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error().Err(err).Msgf("Upgrade error: %v", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}

	h.clientsMu.Lock()
	h.clients[client] = true
	h.clientsMu.Unlock()

	go h.handleWrite(r.Context(), client)
}

func (h *NotificationHandler) Stop() {
	h.sysLogger.Info().Msgf("Stopping the WS notification handler...")

	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()
	for c := range h.clients {
		close(c.send)
		delete(h.clients, c)
	}

	h.sysLogger.Info().Msgf("The WS notification handler stopped.")
}

func (h *NotificationHandler) handleWrite(ctx context.Context, c *Client) {
	logger := log.Ctx(ctx)

	h.sysLogger.Info().Msgf("Handling connection: %v", c)

	defer func() {
		c.conn.Close()
		h.clientsMu.Lock()
		delete(h.clients, c)
		h.clientsMu.Unlock()
	}()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			logger.Error().Err(err).Msgf("Write error: %v", err)
			break
		}
	}
}

func (h *NotificationHandler) broadcastLoop(stop <-chan struct{}) {
	h.sysLogger.Info().Msgf("Broadcast loop started")
	if stop == nil {
		h.sysLogger.Error().Msgf("Broadcast loop stopped.")
	}
	for {
		select {
		case <-stop:
			h.sysLogger.Error().Msgf("Broadcast loop stopped.")
			return
		default:
			notification := Notification{
				Type: "notification",
				Data: NotificationData{
					ID:    25,
					Title: "Премьера фильма:",
					Text:  "Легенда об Очи",
					Date:  "2025-06-08T00:00:00Z",
				},
			}
			msg, _ := json.Marshal(notification)

			h.clientsMu.Lock()
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					// If send buffer is full, drop the client
					close(c.send)
					delete(h.clients, c)
				}
			}
			h.clientsMu.Unlock()

			time.Sleep(60 * time.Second)
		}
	}
}

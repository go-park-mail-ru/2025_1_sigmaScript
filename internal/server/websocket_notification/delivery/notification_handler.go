package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	notifyRefreshCycleTimeInHours = 6
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
	GetMainPageCollections(ctx context.Context) (mocks.Collections, error)
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
				return true
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
		logger.Error().Err(err).Msgf("Upgrade ws connection error: %v", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}

	h.clientsMu.Lock()
	h.clients[client] = true
	h.clientsMu.Unlock()

	go h.handleWrite(r.Context(), client)

	// check for incoming movies once
	errCheckForReleasedMovies := h.checkForReleasedMovies()
	if errCheckForReleasedMovies != nil {
		h.sysLogger.Error().Err(errCheckForReleasedMovies).Msgf("Error happened while checking for released movies: %v", errCheckForReleasedMovies)
	}
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

	defer func() {
		c.conn.Close()
		h.clientsMu.Lock()
		delete(h.clients, c)
		h.clientsMu.Unlock()
	}()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			logger.Error().Err(err).Msgf("Write data to ws connection error: %v", err)
			break
		}
		time.Sleep(2 * time.Second)
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
			errCheckForReleasedMovies := h.checkForReleasedMovies()
			if errCheckForReleasedMovies != nil {
				h.sysLogger.Error().Err(errCheckForReleasedMovies).Msgf("Error happened while checking for released movies: %v", errCheckForReleasedMovies)
				time.Sleep(notifyRefreshCycleTimeInHours * time.Second)
				continue
			}

			time.Sleep(notifyRefreshCycleTimeInHours * time.Hour)
		}
	}
}

func (h *NotificationHandler) checkForReleasedMovies() error {
	collections, err := h.searchService.GetMainPageCollections(context.Background())
	if err != nil {
		return errors.Wrap(err, "Failed to get main page collections")
	}

	now := time.Now().UTC()
	today := now
	tomorrow := now.Add(24 * time.Hour)

	calendarMovieReleases, ok := collections["calendar"]
	if !ok {
		return errors.Wrap(err, "No 'calendar' named collection found")
	}

	for _, movie := range calendarMovieReleases {
		releaseTime, err := time.Parse(time.RFC3339, movie.ReleaseDate)
		if err != nil {
			return errors.Wrapf(err, "Invalid release date format for movie ID %d: %s", movie.ID, movie.ReleaseDate)
		}

		sameDay := releaseTime.Year() == today.Year() &&
			releaseTime.Month() == today.Month() &&
			releaseTime.Day() == today.Day()

		nextDay := releaseTime.Year() == tomorrow.Year() &&
			releaseTime.Month() == tomorrow.Month() &&
			releaseTime.Day() == tomorrow.Day()

		if sameDay || nextDay {
			notification := Notification{
				Type: "notification",
				Data: NotificationData{
					ID:    movie.ID,
					Title: "Премьера фильма:",
					Text:  movie.Title,
					Date:  movie.ReleaseDate,
				},
			}
			msg, errMarshal := json.Marshal(notification)
			if errMarshal != nil {
				return errors.Wrapf(err, "Error happened while marshalling release movie: %v", errMarshal)
			}

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
		}
	}
	return nil
}

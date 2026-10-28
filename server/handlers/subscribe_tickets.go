package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
)

// wsWriteTimeout is how long one write to a client may take before the connection is dropped.
const wsWriteTimeout = 2 * time.Second

func (h *Handler) SubscribeTickets(c *gin.Context) {
	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:*", "127.0.0.1:*"}, // allow local dev pages on any port
	})
	if err != nil {
		h.logger.Error("accept websocket", slog.Any("error", err)) // Accept already wrote the HTTP error

		return
	}
	defer func() { _ = conn.CloseNow() }()

	sub := h.hub.Subscribe()
	defer h.hub.Unsubscribe(sub)

	ctx := conn.CloseRead(c.Request.Context())

	for {
		select {
		case msg := <-sub.Messages():
			data, err := json.Marshal(msg)
			if err != nil {
				h.logger.Error("marshal message to json", slog.Any("error", err))
				continue
			}

			wctx, cancel := context.WithTimeout(ctx, wsWriteTimeout)
			err = conn.Write(wctx, websocket.MessageText, data)

			cancel()

			if err != nil {
				h.logger.Error("send message with websocket", slog.Any("error", err))
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

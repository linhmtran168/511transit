package handler

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/httplog"
	"github.com/linhmtran168/511transit/internal/data"
	"github.com/linhmtran168/511transit/internal/models"
	"github.com/rs/zerolog"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WebSocketHandler struct {
	repository    data.DataRepository
	acceptOptions *websocket.AcceptOptions
}

func NewWebSocketHandler(repository data.DataRepository) *WebSocketHandler {
	acceptOptions := &websocket.AcceptOptions{
		Subprotocols: []string{"realtime-transit"},
	}

	// If not production, add cross origin settings
	if os.Getenv("ENV") != "prod" {
		acceptOptions.OriginPatterns = []string{"localhost*", "127.0.0.1*"}
	}

	return &WebSocketHandler{
		repository:    repository,
		acceptOptions: acceptOptions,
	}
}

func (h *WebSocketHandler) ConnectHandler(w http.ResponseWriter, r *http.Request) {
	httpLogger := httplog.LogEntry(r.Context())
	// https://github.com/nhooyr/websocket/issues/218
	// https://github.com/gorilla/websocket/issues/731
	// Safari doesn't support Compression yet...
	if strings.Contains(r.UserAgent(), "Safari") {
		h.acceptOptions.CompressionMode = websocket.CompressionDisabled
	}
	conn, err := websocket.Accept(w, r, h.acceptOptions)

	if err != nil {
		httpLogger.Error().Err(err).Msg("failed to accept websocket connection")
		return
	}
	defer conn.Close(websocket.StatusInternalError, "internal error happened")

	httpLogger.Info().Msgf("%s", conn.Subprotocol())
	if conn.Subprotocol() != "realtime-transit" {
		conn.Close(websocket.StatusPolicyViolation, "client must speak the correct subprotocol")
		return
	}

	for {
		err := h.processMessage(r.Context(), conn, httpLogger)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
			websocket.CloseStatus(err) == websocket.StatusGoingAway {
			return
		}

		if err != nil {
			httpLogger.Error().Err(err).Msg("Error while handling websocket connection")
			return
		}
	}
}

func (h *WebSocketHandler) processMessage(ctx context.Context, conn *websocket.Conn, httpLogger zerolog.Logger) error {
	// Keep the connection alive until the limit context of 30 minutes
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	var params models.RequestInput
	err := wsjson.Read(ctx, conn, &params)
	if err != nil {
		return err
	}

	switch params.RequestType {
	case "operators":
		operators, err := h.repository.GetOperators()
		if err != nil {
			return err
		}

		err = wsjson.Write(ctx, conn, models.OperatorsResponse{ResponseType: params.RequestType, Data: operators})
		if err != nil {
			return err
		}
	case "tripUpdates":
		operatorID, _ := params.Data["operatorId"].(string)
		tripUpdates, err := h.repository.GetTripUpdates(operatorID)
		if err != nil {
			return err
		}

		tripUpdateData := models.TripUpdateData{OperatorID: operatorID, TripUpdates: tripUpdates}
		response := models.TripUpdatesResponse{ResponseType: params.RequestType, Data: tripUpdateData}
		err = wsjson.Write(ctx, conn, response)
		if err != nil {
			return err
		}
	default:
		httpLogger.Warn().Msgf("unknown request type received: %s", params.RequestType)
	}

	return nil
}

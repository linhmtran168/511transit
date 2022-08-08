package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/httplog"
	"github.com/linhmtran168/511transit/internal/data"
	"github.com/linhmtran168/511transit/internal/models"
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

	// If local, add cross origin settings
	if os.Getenv("ENV") == "" {
		acceptOptions.OriginPatterns = []string{"localhost*", "127.0.0.1*"}
	}

	return &WebSocketHandler{
		repository:    repository,
		acceptOptions: acceptOptions,
	}
}

func (h *WebSocketHandler) ConnectHandler(w http.ResponseWriter, r *http.Request) {
	httpLogger := httplog.LogEntry(r.Context())
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
		err := h.processMessage(r.Context(), conn)
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

func (h *WebSocketHandler) processMessage(ctx context.Context, conn *websocket.Conn) error {
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
		operatorID, _ := params.Data["operatorId"]
		tripUpdate, err := h.repository.GetTripUpdates(fmt.Sprintf("%s", operatorID))
		if err != nil {
			return err
		}

		err = wsjson.Write(ctx, conn, models.TripUpdatesResponse{ResponseType: params.RequestType, Data: tripUpdate})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown request type: %s", params.RequestType)
	}

	return nil
}

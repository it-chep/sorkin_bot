package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
	"sorkin_bot/internal/controller/dto/api/webhook_events"
)

type Controller struct {
	logger              *slog.Logger
	notificationService NotificationService
}

func (c Controller) CancelAppointmentWebhook(ctx *gin.Context) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.logger.Error(fmt.Sprintf("%s", err))
		}
	}(ctx.Request.Body)

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error reading request body: %s", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var requestDTO []webhook_events.AppointmentDTO
	if err = json.Unmarshal(body, &requestDTO); err != nil {
		c.logger.Error(fmt.Sprintf("Error formatting JSON: %s", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if len(requestDTO) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	appointment := requestDTO[0].ToDomain()
	err = c.notificationService.NotifyCancelAppointment(ctx, appointment)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while sending appointment webhook"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func NewController(
	logger *slog.Logger,
	notificationService NotificationService,
) *Controller {
	return &Controller{
		notificationService: notificationService,
		logger:              logger,
	}
}

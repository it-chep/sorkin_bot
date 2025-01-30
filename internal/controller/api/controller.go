package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sorkin_bot/internal/controller/dto/api/webhook_events"
	"strconv"
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

	request, err := parseFormEncodedBody(string(body))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing webhook request"})
		return
	}

	appointment := request.Data.ToDomain()
	err = c.notificationService.NotifyCancelAppointment(ctx, appointment)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while sending appointment webhook"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (c Controller) CreateAppointmentWebhook(ctx *gin.Context) {
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

	request, err := parseFormEncodedBody(string(body))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing webhook request"})
		return
	}

	appointment := request.Data.ToDomain()
	err = c.notificationService.NotifyCreateAppointment(ctx, appointment)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while sending appointment webhook"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func parseFormEncodedBody(body string) (*webhook_events.AppointmentRequest, error) {
	decoded, err := url.ParseQuery(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body: %w", err)
	}

	data := webhook_events.AppointmentDTO{}
	request := &webhook_events.AppointmentRequest{}

	request.Event = decoded.Get("event")
	request.Date = decoded.Get("date")

	data.Id, _ = strconv.Atoi(decoded.Get("data[id]"))
	data.TimeStart = decoded.Get("data[time_start]")
	data.TimeEnd = decoded.Get("data[time_end]")
	data.ClinicId, _ = strconv.Atoi(decoded.Get("data[clinic_id]"))
	data.Clinic = decoded.Get("data[clinic]")
	data.DoctorId, _ = strconv.Atoi(decoded.Get("data[doctor_id]"))
	data.Doctor = decoded.Get("data[doctor]")
	data.PatientId, _ = strconv.Atoi(decoded.Get("data[patient_id]"))
	data.PatientName = decoded.Get("data[patient_name]")
	data.PatientBirthDate = decoded.Get("data[patient_birth_date]")
	data.PatientGender = decoded.Get("data[patient_gender]")
	data.DateCreated = decoded.Get("data[date_created]")
	data.DateUpdated = decoded.Get("data[date_updated]")
	data.Status = decoded.Get("data[status]")
	data.StatusId, _ = strconv.Atoi(decoded.Get("data[status_id]"))
	data.PatientPhone = decoded.Get("data[patient_phone]")

	request.Data = data

	return request, nil
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

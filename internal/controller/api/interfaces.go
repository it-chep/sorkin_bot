package api

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
)

type NotificationService interface {
	NotifyCancelAppointment(ctx context.Context, appointment appointment.Appointment) error
	NotifyCreateAppointment(ctx context.Context, appointment appointment.Appointment) error
}

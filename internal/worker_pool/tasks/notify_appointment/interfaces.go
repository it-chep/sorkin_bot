package notify_appointment

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
)

type appointmentService interface {
	GetAppointmentsForNotifying(ctx context.Context) ([]appointment.Appointment, error)
}

type notificationService interface {
	NotifySoonAppointment(ctx context.Context, appointment appointment.Appointment) error
}

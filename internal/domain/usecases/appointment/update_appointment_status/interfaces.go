package update_appointment_status

import "context"

type writeRepo interface {
	UpdateStatusDraftAppointment(ctx context.Context, tgId int64, appointmentId int) (err error)
}

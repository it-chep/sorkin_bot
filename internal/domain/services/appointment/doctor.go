package appointment

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
)

func (as *AppointmentService) GetDoctors(ctx context.Context, specialityId int) (doctors []appointment.Doctor) {
	_ = "sorkin_bot.internal.domain.services.appointment.doctor.GetDoctors"

	doctors = as.misAdapter.GetDoctors(ctx, specialityId)

	return doctors
}

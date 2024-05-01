package appointment

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
)

func (as *AppointmentService) GetDoctors(ctx context.Context, specialityId int) (doctors []appointment.Doctor) {
	op := "sorkin_bot.internal.domain.services.appointment.doctor.GetDoctors"

	err, doctors := as.mis.GetDoctors(ctx, specialityId)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return doctors
	}

	return doctors
}

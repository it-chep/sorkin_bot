package appointment

import (
	"context"
	"fmt"
	entity "sorkin_bot/internal/domain/entity/user"
)

func (as AppointmentService) GetPatient(ctx context.Context, user entity.User) (result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.user.GetPatient"
	err := as.misAdapter.GetPatientById(ctx, *user.GetPatientId())
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return false
	}
	return true
}

func (as AppointmentService) CreatePatient(ctx context.Context, user entity.User) (result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.user.CreatePatient"

	err, patientId := as.misAdapter.CreatePatient(ctx, user)
	err = as.userService.UpdatePatientId(ctx, user, patientId)

	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return false
	}

	return true
}

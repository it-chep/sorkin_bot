package appointment

import (
	"context"
	"fmt"
	"sorkin_bot/internal/clients/gateways/dto"
	entity "sorkin_bot/internal/domain/entity/user"
)

func (as *AppointmentService) GetPatient(ctx context.Context, user entity.User) (result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.user.GetPatient"
	if user.GetPatientId() == nil {
		return false
	}
	_, err := as.misAdapter.GetPatientById(ctx, *user.GetPatientId())
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return false
	}
	return true
}

func (as *AppointmentService) CreatePatient(ctx context.Context, user entity.User) (result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.user.CreatePatient"
	userDTO := dto.PatientDTO{
		LastName:  *user.GetLastName(),
		FirstName: user.GetFirstName(),
		ThirdName: user.GetThirdName(),
		BirthDate: *user.GetBirthDate(),
		Phone:     *user.GetPhone(),
	}
	patientId, err := as.misAdapter.CreatePatient(ctx, userDTO)
	err = as.userService.UpdatePatientId(ctx, user, patientId)

	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return false
	}

	return true
}

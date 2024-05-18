package appointment

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
	"sorkin_bot/pkg/utils"
)

func (as *AppointmentService) GetDoctors(ctx context.Context, tgId int64, offset int, specialityId *int) (doctorsMap map[int]string) {
	_ = "sorkin_bot.internal.domain.services.appointment.doctor.GetDoctors"
	if specialityId == nil {
		draftAppointment, err := as.GetDraftAppointment(ctx, tgId)
		if err != nil {
			return nil
		}
		specialityIdValue := draftAppointment.GetSpecialityId()
		specialityId = &specialityIdValue
	}
	doctors := as.misAdapter.GetDoctors(ctx, *specialityId)
	doctorsMap = as.getDoctorsMap(doctors)
	doctorsMap = utils.IntMapWithOffset(utils.SortedIntMap(doctorsMap), offset)

	return doctorsMap
}

func (as *AppointmentService) getDoctorsMap(doctors []appointment.Doctor) (doctorsMap map[int]string) {
	doctorsMap = make(map[int]string)
	for _, doctor := range doctors {
		doctorsMap[doctor.GetID()] = doctor.GetName()
	}
	return doctorsMap
}

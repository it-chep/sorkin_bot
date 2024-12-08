package appointment

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/pkg/utils"
)

func (as *AppointmentService) GetDoctorsBySpecialityId(ctx context.Context, tgId int64, offset int, specialityId *int) (doctorsMap map[int]string) {
	_ = "sorkin_bot.internal.domain.services.appointment.doctor.GetDoctorsBySpecialityId"
	if specialityId == nil {
		draftAppointment, err := as.GetDraftAppointment(ctx, tgId)
		if err != nil {
			return nil
		}
		specialityIdValue := draftAppointment.GetSpecialityId()
		specialityId = specialityIdValue
	}

	doctors := as.misAdapter.GetDoctorsBySpecialityId(ctx, *specialityId)
	doctorsMap = as.getDoctorsMap(doctors)
	doctorsMap = utils.IntMapWithOffset(utils.SortedIntMap(doctorsMap), offset)

	return doctorsMap
}

func (as *AppointmentService) GetDoctors(ctx context.Context, tgId int64, offset int) (doctorsMap map[int]string) {
	_ = "sorkin_bot.internal.domain.services.appointment.doctor.GetDoctors"
	var doctors []appointment.Doctor
	draftAppointment, err := as.readDraftAppointmentRepo.GetUserDraftAppointment(ctx, tgId)
	if err != nil {
		return nil
	}
	if draftAppointment.GetAppointmentType() == nil {
		return nil
	}
	switch *draftAppointment.GetAppointmentType() {
	case appointment.OnlineAppointment:
		doctors, err = as.misAdapter.GetDoctors(ctx, false, true, false)
	case appointment.ClinicAppointment:
		doctors, err = as.misAdapter.GetDoctors(ctx, false, false, true)
	case appointment.HomeAppointment:
		doctors, err = as.misAdapter.GetDoctors(ctx, true, false, false)
	}
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

func (as *AppointmentService) GetDoctorInfo(ctx context.Context, user entity.User, doctorId int) (doctorEntity appointment.Doctor, err error) {
	var translations []appointment.TranslationEntity
	doctor := as.misAdapter.GetDoctorInfo(ctx, doctorId)
	ids := doctor.GetSecondProfessions()
	translationString := ""
	if len(ids) > 0 {
		translations, err = as.readMessageRepo.GetManyTranslationsByIds(ctx, ids)
		if err != nil {
			return appointment.Doctor{}, err
		}
		for _, translation := range translations {
			translationString = fmt.Sprintf(
				"%s %s", as.GetTranslationString(*user.GetLanguageCode(), translation), translationString,
			)
		}
		doctorEntity = doctor.SetDoctorInfo(translationString)
		return doctorEntity, nil
	} else {
		return appointment.Doctor{}, fmt.Errorf("doctor not found")
	}
}

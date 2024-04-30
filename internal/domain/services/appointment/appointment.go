package appointment

import (
	"context"
	"fmt"
	"log/slog"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/services/user"
)

type AppointmentService struct {
	Mis         Appointment
	userService user.UserService
	readRepo    ReadRepo
	logger      *slog.Logger
}

func NewAppointmentService(mis Appointment, readRepo ReadRepo, logger *slog.Logger, userService user.UserService) AppointmentService {
	return AppointmentService{
		Mis:         mis,
		userService: userService,
		readRepo:    readRepo,
		logger:      logger,
	}
}

func (as *AppointmentService) GetSchedules(ctx context.Context, doctorId int) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.GetPatient"
	err, _ := as.Mis.GetSchedules(ctx, doctorId)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return
	}
}

func (as *AppointmentService) GetPatient(ctx context.Context, user entity.User) (result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.GetPatient"
	err := as.Mis.GetPatientById(ctx, user.GetPatientId())
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return false
	}
	return true
}

func (as *AppointmentService) CreatePatient(ctx context.Context, user entity.User) (result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.CreatePatient"

	err, patientId := as.Mis.CreatePatient(ctx, user)
	err = as.userService.UpdatePatientId(ctx, user, patientId)

	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return false
	}

	return true
}

func (as *AppointmentService) GetTranslatedSpecialities(
	ctx context.Context,
	user entity.User,
	specialities []appointment.Speciality,
) (translatedSpecialities map[int]string, err error) {
	var translatedSpeciality string
	op := "sorkin_bot.internal.domain.services.appointment.appointment.GetTranslatedSpecialities"
	translations, err := as.readRepo.GetTranslationsBySlug(ctx, "doctor")
	translatedSpecialities = make(map[int]string)

	if err != nil {
		return translatedSpecialities, err
	}
	langCode := user.GetLanguageCode()

	for _, speciality := range specialities {
		translationEntity, ok := translations[speciality.GetDoctorName()]

		if !ok {
			as.logger.Error(fmt.Sprintf("untranslated speciality: %s, please translate this in priority. Place %s", speciality.GetDoctorName(), op))
			translatedSpeciality = speciality.GetDoctorName()
		}

		switch langCode {
		case "RU":
			translatedSpeciality = translationEntity.GetRuText()
		case "EN":
			translatedSpeciality = translationEntity.GetEngText()
		case "PT":
			translatedSpeciality = translationEntity.GetPtBrText()
		}

		translatedSpecialities[speciality.GetId()] = translatedSpeciality
	}
	return translatedSpecialities, err
}

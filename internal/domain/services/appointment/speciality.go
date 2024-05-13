package appointment

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

func (as *AppointmentService) GetSpecialities(ctx context.Context) (specialities []appointment.Speciality, err error) {
	specialities = as.misAdapter.GetSpecialities(ctx)
	if err != nil {
		return nil, err
	}
	return specialities, err
}

func (as *AppointmentService) GetTranslatedSpecialities(
	ctx context.Context,
	user entity.User,
	specialities []appointment.Speciality,
) (translatedSpecialities map[int]string, unTranslatedSpecialities []string, err error) {
	var translatedSpeciality string
	op := "sorkin_bot.internal.domain.services.appointment.speciality.GetTranslatedSpecialities"
	translations, err := as.readRepo.GetTranslationsBySlug(ctx, "doctor")
	translatedSpecialities = make(map[int]string)
	unTranslatedSpecialities = make([]string, 0)

	if err != nil {
		return translatedSpecialities, nil, err
	}
	langCode := user.GetLanguageCode()

	for _, speciality := range specialities {
		translationEntity, ok := translations[speciality.GetDoctorName()]

		if !ok {
			as.logger.Error(fmt.Sprintf("untranslated speciality: %s, please translate this in priority. Place %s", speciality.GetDoctorName(), op))
			unTranslatedSpecialities = append(unTranslatedSpecialities, speciality.GetDoctorName())
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
	return translatedSpecialities, unTranslatedSpecialities, err
}

package appointment

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

func (as AppointmentService) GetTranslatedSpecialities(
	ctx context.Context,
	user entity.User,
	specialities []appointment.Speciality,
) (translatedSpecialities map[int]string, err error) {
	var translatedSpeciality string
	op := "sorkin_bot.internal.domain.services.appointment.speciality.GetTranslatedSpecialities"
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

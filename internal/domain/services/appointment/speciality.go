package appointment

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/pkg/utils"
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
	offset int,
) (translatedSpecialities map[int]string, unTranslatedSpecialities []string, err error) {
	var translatedSpeciality string
	op := "sorkin_bot.internal.domain.services.appointment.speciality.GetTranslatedSpecialities"
	translations, err := as.readMessageRepo.GetTranslationsBySlug(ctx, "doctor")
	translatedSpecialities = make(map[int]string)
	unTranslatedSpecialities = make([]string, 0)

	if err != nil {
		return translatedSpecialities, nil, err
	}
	langCode := *user.GetLanguageCode()

	for _, speciality := range specialities {
		translationEntity, ok := translations[speciality.GetDoctorName()]

		if !ok && speciality.GetDoctorName() != "" {
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

		if translatedSpeciality == "" || !translations[speciality.GetDoctorName()].GetUses() {
			continue
		}

		translatedSpecialities[speciality.GetId()] = translatedSpeciality
	}
	translatedSpecialities = utils.IntMapWithOffset(utils.SortedIntMap(translatedSpecialities), offset)
	return translatedSpecialities, unTranslatedSpecialities, err
}

func (as *AppointmentService) TranslateSpecialityByID(ctx context.Context, user entity.User, specialityId int) (translatedSpeciality string, err error) {

	translationEntity, err := as.readMessageRepo.GetTranslationsBySourceId(ctx, specialityId)
	if err != nil {
		return "", err
	}

	langCode := *user.GetLanguageCode()
	switch langCode {
	case "RU":
		translatedSpeciality = translationEntity.GetRuText()
	case "EN":
		translatedSpeciality = translationEntity.GetEngText()
	case "PT":
		translatedSpeciality = translationEntity.GetPtBrText()
	}

	return translatedSpeciality, nil
}

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
	return specialities, nil
}

func (as *AppointmentService) GetTranslatedSpecialities(
	ctx context.Context,
	user entity.User,
	specialities []appointment.Speciality,
	offset int,
) (translatedSpecialities map[int]string, unTranslatedSpecialities []string, err error) {
	var translatedSpeciality string
	op := "sorkin_bot.internal.domain.services.appointment.speciality.GetTranslatedSpecialities"
	translations, err := as.readMessageRepo.GetTranslationsBySlugKeySlug(ctx, "Дополнительно")
	translatedSpecialities = make(map[int]string)
	unTranslatedSpecialities = make([]string, 0)

	if err != nil {
		return translatedSpecialities, nil, err
	}
	langCode := *user.GetLanguageCode()

	for _, speciality := range specialities {
		translationEntity, ok := translations[speciality.GetName()]

		if !ok && speciality.GetName() != "" {
			as.logger.Error(fmt.Sprintf("untranslated speciality: %s, please translate this in priority. Place %s", speciality.GetDoctorName(), op))
			unTranslatedSpecialities = append(unTranslatedSpecialities, speciality.GetDoctorName())
		}

		translatedSpeciality = as.GetTranslationString(langCode, translationEntity)

		if translatedSpeciality == "" || !translations[speciality.GetName()].GetUses() {
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
	translatedSpeciality = as.GetTranslationString(langCode, translationEntity)

	return translatedSpeciality, nil
}

func (as *AppointmentService) TranslateManyByIds(ctx context.Context, user entity.User, ids []int) (translatedSpecialities map[int]string, err error) {
	translations, err := as.readMessageRepo.GetManyTranslationsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	translatedSpecialities = make(map[int]string)

	for _, translationEntity := range translations {

		langCode := *user.GetLanguageCode()
		translatedSpeciality := as.GetTranslationString(langCode, translationEntity)

		if translationEntity.GetSourceId() != nil {
			translatedSpecialities[*translationEntity.GetSourceId()] = translatedSpeciality
		}

	}
	return translatedSpecialities, err
}

func (as *AppointmentService) GetTranslationString(langCode string, translationEntity appointment.TranslationEntity) (translatedString string) {

	switch langCode {
	case "RU":
		translatedString = translationEntity.GetRuText()
	case "EN":
		translatedString = translationEntity.GetEngText()
	case "PT":
		translatedString = translationEntity.GetPtBrText()
	}

	return translatedString
}

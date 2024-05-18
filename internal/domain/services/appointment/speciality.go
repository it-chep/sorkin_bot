package appointment

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"sort"
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
	translations, err := as.readRepo.GetTranslationsBySlug(ctx, "doctor")
	translatedSpecialities = make(map[int]string)
	unTranslatedSpecialities = make([]string, 0)

	if err != nil {
		return translatedSpecialities, nil, err
	}
	langCode := user.GetLanguageCode()

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

		if translatedSpeciality == "" {
			continue
		}

		translatedSpecialities[speciality.GetId()] = translatedSpeciality
	}
	translatedSpecialities = as.specialitiesWithOffset(as.sortedSpeciality(translatedSpecialities), offset)
	return translatedSpecialities, unTranslatedSpecialities, err
}

func (as *AppointmentService) sortedSpeciality(m map[int]string) map[int]string {
	sortedKeys := make([]int, 0, len(m))
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)

	sortedMap := make(map[int]string)
	for _, k := range sortedKeys {
		sortedMap[k] = m[k]
	}
	return sortedMap
}

func (as *AppointmentService) specialitiesWithOffset(allSpecialities map[int]string, offset int) map[int]string {
	sortedKeys := make([]int, 0, len(allSpecialities))
	for k := range allSpecialities {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)
	offsetMap := make(map[int]string)
	if offset < len(sortedKeys) {
		for _, k := range sortedKeys[offset:] {
			offsetMap[k] = allSpecialities[k]
		}
	}
	return offsetMap
}

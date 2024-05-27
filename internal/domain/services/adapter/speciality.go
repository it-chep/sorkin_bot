package adapter

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
	"time"
)

func (a *AppointmentServiceAdapter) GetSpecialities(ctx context.Context) (specialities []appointment.Speciality) {
	cachedSpecialities, ok := a.cache.Get("specialities")
	if !ok {
		specialitiesDTO, err := a.gateway.GetSpecialities(ctx)
		if err != nil {
			return
		}

		for _, specialityDTO := range specialitiesDTO {
			specialities = append(specialities, specialityDTO.ToDomain())
		}
		a.cache.Set("specialities", specialities, 12*time.Hour)
		return specialities
	}
	return cachedSpecialities.([]appointment.Speciality)
}

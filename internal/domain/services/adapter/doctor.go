package adapter

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
	"time"
)

func (a *AppointmentServiceAdapter) GetDoctors(ctx context.Context, specialityId int) (doctors []appointment.Doctor) {
	//cachedDoctors, ok := a.cache.Get("doctors")
	//if !ok {
	doctorsDTO, err := a.gateway.GetDoctors(ctx, specialityId)
	if err != nil {
		return
	}

	for _, doctorDTO := range doctorsDTO {
		doctors = append(doctors, doctorDTO.ToDomain())
	}
	a.cache.Set("doctors", doctors, 12*time.Hour)
	return doctors
	//}
	//return cachedDoctors.([]appointment.Doctor)
}

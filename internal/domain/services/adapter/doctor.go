package adapter

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
	"time"
)

func (a *AppointmentServiceAdapter) GetDoctorsBySpecialityId(ctx context.Context, specialityId int) (doctors []appointment.Doctor) {
	doctorsDTO, err := a.gateway.GetDoctorsBySpecialityId(ctx, specialityId)
	if err != nil {
		return
	}

	for _, doctorDTO := range doctorsDTO {
		doctors = append(doctors, doctorDTO.ToDomain())
	}
	a.cache.Set("doctors", doctors, 12*time.Hour)
	return doctors
}

func (a *AppointmentServiceAdapter) GetDoctors(ctx context.Context, homeVisit, onlineAppointment, clinicAppointment bool) (doctors []appointment.Doctor, err error) {
	doctorsDTO, err := a.gateway.GetDoctors(ctx, homeVisit, onlineAppointment, clinicAppointment)
	if err != nil {
		return nil, err
	}

	for _, doctorDTO := range doctorsDTO {
		doctors = append(doctors, doctorDTO.ToDomain())
	}

	// todo add cache
	return doctors, nil
}

func (a *AppointmentServiceAdapter) GetDoctorInfo(ctx context.Context, doctorId int) (doctor appointment.Doctor) {
	cachedDoctors, ok := a.cache.Get(fmt.Sprintf("%d_doctors", doctorId))
	if !ok {
		doctorDTO, err := a.gateway.GetDoctorInfo(ctx, doctorId)
		if err != nil {
			return
		}
		doctor = doctorDTO.ToDomain()
		a.cache.Set(fmt.Sprintf("%d_doctors", doctorId), doctor, 12*time.Hour)
		return doctor
	}
	return cachedDoctors.(appointment.Doctor)

}

package dto

import "sorkin_bot/internal/domain/entity/appointment"

type SpecialityDTO struct {
	ID         int
	Name       string
	DoctorName string
	IsDeleted  bool
}

func NewSpecialityDTO(id int, name, doctorName string, isDeleted bool) SpecialityDTO {
	return SpecialityDTO{
		ID:         id,
		Name:       name,
		DoctorName: doctorName,
		IsDeleted:  isDeleted,
	}
}

func (s SpecialityDTO) ToDomain() appointment.Speciality {
	return appointment.NewSpeciality(
		s.ID, s.Name, s.DoctorName, s.IsDeleted,
	)
}

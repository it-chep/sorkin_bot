package dto

import "sorkin_bot/internal/domain/entity/appointment"

type DoctorDTO struct {
	id                     int
	name                   string
	phone                  string
	email                  string
	professionTitles       string
	secondProfessionTitles string
	isDeleted              bool
}

func NewDoctorDTO(id int, name, phone, email, professionTitles, secondProfessionTitles string, isDeleted bool) DoctorDTO {
	return DoctorDTO{
		id:                     id,
		name:                   name,
		phone:                  phone,
		email:                  email,
		professionTitles:       professionTitles,
		secondProfessionTitles: secondProfessionTitles,
		isDeleted:              isDeleted,
	}
}

func (d DoctorDTO) ToDomain() appointment.Doctor {
	return appointment.NewDoctor(
		d.id, d.name, d.phone, d.email, d.professionTitles, d.secondProfessionTitles, d.isDeleted,
	)
}

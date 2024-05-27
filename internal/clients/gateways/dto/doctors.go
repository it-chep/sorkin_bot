package dto

import "sorkin_bot/internal/domain/entity/appointment"

type DoctorDTO struct {
	id                     int
	name                   string
	phone                  string
	email                  string
	professionTitles       string
	secondProfessions      []int
	secondProfessionTitles string
	isDeleted              bool
}

func NewDoctorDTO(id int, name, phone, email, professionTitles, secondProfessionTitles string, isDeleted bool, secondProfessions []int) DoctorDTO {
	return DoctorDTO{
		id:                     id,
		name:                   name,
		phone:                  phone,
		email:                  email,
		professionTitles:       professionTitles,
		secondProfessions:      secondProfessions,
		secondProfessionTitles: secondProfessionTitles,
		isDeleted:              isDeleted,
	}
}

func (d DoctorDTO) ToDomain() appointment.Doctor {
	return appointment.NewDoctor(
		d.id, d.name, d.phone, d.email, d.professionTitles, d.secondProfessionTitles, d.isDeleted, d.secondProfessions,
	)
}

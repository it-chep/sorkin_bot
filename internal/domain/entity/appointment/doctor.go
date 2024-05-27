package appointment

type Doctor struct {
	id                     int
	name                   string
	phone                  string
	email                  string
	professionTitles       string
	secondProfessions      []int
	secondProfessionTitles string
	isDeleted              bool
	doctorInfo             string
}

func NewDoctor(id int, name, phone, email, professionTitles, secondProfessionTitles string, isDeleted bool, secondProfessions []int) Doctor {
	return Doctor{
		id:                     id,
		name:                   name,
		phone:                  phone,
		email:                  email,
		professionTitles:       professionTitles,
		secondProfessions:      secondProfessions,
		secondProfessionTitles: secondProfessionTitles,
		isDeleted:              isDeleted,
		doctorInfo:             "",
	}
}

func (d *Doctor) GetID() int {
	return d.id
}

func (d *Doctor) GetName() string {
	return d.name
}

func (d *Doctor) GetSecondProfessions() []int {
	return d.secondProfessions
}

func (d *Doctor) SetDoctorInfo(doctorInfo string) Doctor {
	d.doctorInfo = doctorInfo
	return *d
}

func (d *Doctor) GetInfo() string {
	return d.doctorInfo
}

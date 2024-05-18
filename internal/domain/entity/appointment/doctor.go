package appointment

type Doctor struct {
	id                     int
	name                   string
	phone                  string
	email                  string
	professionTitles       string
	secondProfessionTitles string
	isDeleted              bool
}

func NewDoctor(id int, name, phone, email, professionTitles, secondProfessionTitles string, isDeleted bool) Doctor {
	return Doctor{
		id:                     id,
		name:                   name,
		phone:                  phone,
		email:                  email,
		professionTitles:       professionTitles,
		secondProfessionTitles: secondProfessionTitles,
		isDeleted:              isDeleted,
	}
}

func (d *Doctor) GetID() int {
	return d.id
}

func (d *Doctor) GetName() string {
	return d.name
}

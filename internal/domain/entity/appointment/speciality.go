package appointment

type Speciality struct {
	id         int
	name       string
	doctorName string
	isDeleted  bool
}

func NewSpeciality(id int, name, doctorName string, isDeleted bool) Speciality {
	return Speciality{
		id:         id,
		name:       name,
		doctorName: doctorName,
		isDeleted:  isDeleted,
	}
}

func (s Speciality) GetDoctorName() string {
	return s.doctorName
}

func (s Speciality) GetId() int {
	return s.id
}

func (s Speciality) GetName() string { return s.name }

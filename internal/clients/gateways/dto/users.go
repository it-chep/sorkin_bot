package dto

type CreatedPatientDTO struct {
	PatientID int
	Number    int
	Age       string
	Gender    string
	Mobile    string
	Email     string
	PatientDTO
}

func NewCreatePatientDTO() CreatedPatientDTO {
	return CreatedPatientDTO{}
}

func (d CreatedPatientDTO) ToDomain() {
	return
}

type PatientDTO struct {
	LastName  string
	FirstName string
	ThirdName string
	BirthDate string
	Phone     string
	Email     string
}

func NewPatientDTO() PatientDTO {
	return PatientDTO{}
}

func (d PatientDTO) ToDomain() {

}

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

func NewCreatePatientDTO(patientId, number int) CreatedPatientDTO {
	return CreatedPatientDTO{
		PatientID: patientId,
		Number:    number,
	}
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

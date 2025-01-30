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

func NewCreatePatientDTO(patientId, number int, age, gender, mobile, email, phone string) CreatedPatientDTO {
	return CreatedPatientDTO{
		PatientID: patientId,
		Number:    number,
		Age:       age,
		Gender:    gender,
		Mobile:    mobile,
		Email:     email,
		PatientDTO: PatientDTO{
			Phone: phone,
		},
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

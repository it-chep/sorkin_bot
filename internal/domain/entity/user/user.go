package entity

type User struct {
	firstName        string
	lastName         *string
	tgID             int64
	isBot            bool
	registrationTime string
	birthDate        *string
	username         *string
	languageCode     *string
	state            *string
	phone            *string
	patientId        *int
	thirdName        string
}

func NewUser(tgId int64, firstName string, opts ...UserOpt) *User {
	u := &User{
		tgID:      tgId,
		firstName: firstName,
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

func (usr *User) GetBirthDate() *string {
	return usr.birthDate
}

func (usr *User) GetThirdName() string { return usr.thirdName }

func (usr *User) GetPhone() *string {
	return usr.phone
}

func (usr *User) GetRegistrationTime() string {
	return usr.registrationTime
}

func (usr *User) GetPatientId() *int {
	return usr.patientId
}

func (usr *User) GetFirstName() string {
	return usr.firstName
}

func (usr *User) GetLastName() *string {
	return usr.lastName
}

func (usr *User) GetTgId() int64 {
	return usr.tgID
}

func (usr *User) GetUsername() *string {
	return usr.username
}

func (usr *User) GetLanguageCode() *string {
	return usr.languageCode
}

func (usr *User) GetState() *string {
	return usr.state
}

// SetState обновляет текущее состояние пользователя.
func (usr *User) SetState(newState string) {
	usr.state = &newState
}

// SetPatientId обновляет id пользователя в мис
func (usr *User) SetPatientId(patientId *int) {
	usr.patientId = patientId
}

func (usr *User) SetLanguageCode(languageCode string) {
	usr.languageCode = &languageCode
}

type Appointment struct {
}

package entity

type User struct {
	firstName    string
	lastName     string
	tgID         int64
	isBot        bool
	username     string
	languageCode string
	state        string
	phone        string
	patientId    int64
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

func (usr *User) GetPatientId() int64 {
	return usr.patientId
}

func (usr *User) GetFirstName() string {
	return usr.firstName
}

func (usr *User) GetLastName() string {
	return usr.lastName
}

func (usr *User) GetTgId() int64 {
	return usr.tgID
}

func (usr *User) GetUsername() string {
	return usr.username
}

func (usr *User) GetLanguageCode() string {
	return usr.languageCode
}

func (usr *User) GetState() string {
	return usr.state
}

// SetState обновляет текущее состояние пользователя.
func (usr *User) SetState(newState string) {
	usr.state = newState
}

// SetPatientId обновляет id пользователя в мис
func (usr *User) SetPatientId(patientId int64) {
	usr.patientId = patientId
}

func (usr *User) SetLanguageCode(languageCode string) {
	usr.languageCode = languageCode
}

type Appointment struct {
}

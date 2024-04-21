package entity

type User struct {
	firstName    string
	lastName     string
	tgID         int64
	isBot        bool
	username     string
	languageCode string
	state        string
}

func NewUser(tgId int64, firstName string, opts ...UserOpt) *User {
	u := &User{
		tgID: tgId,
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
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

type Appointment struct {
}

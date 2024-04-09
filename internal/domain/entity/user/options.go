package entity

type UserOpt func(usr *User) *User

func WithUsrLanguageCode(languageCode string) UserOpt {
	return func(usr *User) *User {
		usr.languageCode = languageCode
		return usr
	}
}

func WithUsrUsername(username string) UserOpt {
	return func(usr *User) *User {
		usr.username = username
		return usr
	}
}

func WithUsrLastName(lastName string) UserOpt {
	return func(usr *User) *User {
		usr.lastName = lastName
		return usr
	}
}

func WithState(state string) UserOpt {
	return func(usr *User) *User {
		usr.state = state
		return usr
	}
}

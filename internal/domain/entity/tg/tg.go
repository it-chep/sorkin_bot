package tg

import "time"

type MessageLog struct {
	id              int64
	tgMessageId     int64
	systemMessageId int64
	userId          int64
	timestamp       time.Time
}

type MessageCondition struct {
	id int64
}

type Message struct {
	id int64
}

type TgUser struct {
	id               int64
	tgId             int64
	name             string
	surname          string
	username         string
	phone            string
	lastState        string
	registrationDate time.Time
}

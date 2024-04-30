package tg

import "time"

type MessageLog struct {
	tgMessageId int64
	messageText string
	userId      int64
}

func NewMessageLog(tgMessageId, userId int64, messageText string) MessageLog {
	return MessageLog{
		tgMessageId: tgMessageId,
		messageText: messageText,
		userId:      userId,
	}
}

func (ml *MessageLog) GetTgMessageId() int64 {
	return ml.tgMessageId
}

func (ml *MessageLog) GetMessageText() string {
	return ml.messageText
}

func (ml *MessageLog) GetUserTgId() int64 {
	return ml.userId
}

type MessageCondition struct {
	id int64
}

type Message struct {
	id       int64
	ruText   string
	engText  string
	ptBrText string
}

func NewMessage(id int64, ruText, engText, ptBrText string) Message {
	return Message{
		id:       id,
		ruText:   ruText,
		engText:  engText,
		ptBrText: ptBrText,
	}
}

func (m Message) GetRuText() string {
	return m.ruText
}

func (m Message) GetEngText() string {
	return m.engText
}

func (m Message) GetRtBRText() string {
	return m.ptBrText
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

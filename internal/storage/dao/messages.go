package dao

import (
	"sorkin_bot/internal/domain/entity/tg"
	"time"
)

type MessageDAO struct {
	Id       int64  `db:"id"`
	RuText   string `db:"ru_text"`
	EngText  string `db:"eng_text"`
	PtBrText string `db:"pt_br_text"`
}

func (m MessageDAO) ToDomain() tg.Message {
	return tg.NewMessage(
		m.Id, m.RuText, m.EngText, m.PtBrText,
	)
}

type MessageLogDAO struct {
	Id          int64     `db:"id"`
	TgMessageID int64     `db:"tg_message_id"`
	Text        string    `db:"text"`
	UserTgID    int64     `db:"user_tg_id"`
	Time        time.Time `db:"time"`
}

func (m MessageLogDAO) ToDomain() tg.MessageLog {
	return tg.NewMessageLog(m.TgMessageID, m.UserTgID, m.Text)
}

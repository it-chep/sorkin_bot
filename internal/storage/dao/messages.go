package dao

import "sorkin_bot/internal/domain/entity/tg"

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

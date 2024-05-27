package dao

import "sorkin_bot/internal/domain/entity/appointment"

type TranslationDao struct {
	Slug       string `db:"slug"`
	RuText     string `db:"ru_text"`
	EngText    string `db:"eng_text"`
	PtBrText   string `db:"pt_br_text"`
	Uses       bool   `db:"uses"`
	SourceId   *int   `db:"id_in_source_system"`
	Profession string `db:"profession"`
}

func (dao *TranslationDao) ToDomain() appointment.TranslationEntity {
	return appointment.NewTranslationEntity(
		dao.Slug, dao.Profession, dao.RuText, dao.EngText, dao.PtBrText, dao.Uses, dao.SourceId,
	)
}

package appointment

type TranslationEntity struct {
	slug     string
	ruText   string
	engText  string
	ptBrText string
}

func NewTranslationEntity(slug, ruText, engText, ptBrText string) TranslationEntity {
	return TranslationEntity{
		slug:     slug,
		ruText:   ruText,
		engText:  engText,
		ptBrText: ptBrText,
	}
}

func (te TranslationEntity) GetRuText() string {
	return te.ruText
}

func (te TranslationEntity) GetEngText() string {
	return te.engText
}

func (te TranslationEntity) GetPtBrText() string {
	return te.ptBrText
}

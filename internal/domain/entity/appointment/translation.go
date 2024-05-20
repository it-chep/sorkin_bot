package appointment

type TranslationEntity struct {
	slug     string
	ruText   string
	engText  string
	ptBrText string
	uses     bool
	sourceId *int
}

func NewTranslationEntity(slug, ruText, engText, ptBrText string, uses bool, sourceId *int) TranslationEntity {
	return TranslationEntity{
		slug:     slug,
		ruText:   ruText,
		engText:  engText,
		ptBrText: ptBrText,
		uses:     uses,
		sourceId: sourceId,
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

func (te TranslationEntity) GetUses() bool {
	return te.uses
}

func (te TranslationEntity) GetSourceId() *int {
	return te.sourceId
}

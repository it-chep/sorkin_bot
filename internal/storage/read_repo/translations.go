package read_repo

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"log/slog"
	"sorkin_bot/internal/domain/entity/appointment"
	"sorkin_bot/internal/storage/dao"
	"sorkin_bot/pkg/client/postgres"
)

type TranslationStorage struct {
	client postgres.Client
	logger *slog.Logger
}

func NewTranslationRepo(client postgres.Client, logger *slog.Logger) TranslationStorage {
	return TranslationStorage{
		client: client,
		logger: logger,
	}
}

func (tr TranslationStorage) GetTranslationsBySlug(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error) {
	var translationsDao []dao.TranslationDao
	op := "sorkin_bot.internal.storage.read_repo.translations.GetTranslationsBySlug"
	q := `select slug, ru_text, eng_text, pt_br_text, uses from translations where slug like  $1 || '%';`

	err = pgxscan.Select(ctx, tr.client, &translationsDao, q, slug)
	if err != nil {
		tr.logger.Error(fmt.Sprintf("Error while scanning row: %s op: %s", err, op))
		return translations, err
	}

	translations = make(map[string]appointment.TranslationEntity)
	for _, translation := range translationsDao {
		translations[translation.RuText] = translation.ToDomain()
	}
	return translations, nil
}

func (tr TranslationStorage) GetTranslationsBySourceId(ctx context.Context, sourceId int) (translation appointment.TranslationEntity, err error) {
	var translationsDao dao.TranslationDao
	op := "sorkin_bot.internal.storage.read_repo.translations.GetTranslationsBySourceId"
	q := `select slug, ru_text, eng_text, pt_br_text, uses from translations where id_in_source_system = $1 and uses = true;`

	err = pgxscan.Get(ctx, tr.client, &translationsDao, q, sourceId)
	if err != nil {
		tr.logger.Error(fmt.Sprintf("Error while scanning row: %s op: %s", err, op))
		return translation, err
	}
	translation = translationsDao.ToDomain()

	return translation, nil
}

package read_repo

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"log/slog"
	"sorkin_bot/internal/domain/entity/appointment"
	"sorkin_bot/internal/storage/dao"
	"sorkin_bot/pkg/client/postgres"
	"strings"
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

func (tr TranslationStorage) GetTranslationsBySlugKeySlug(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error) {
	var translationsDao []dao.TranslationDao
	op := "sorkin_bot.internal.storage.read_repo.translations.GetTranslationsBySlugKeySlug"
	q := `select slug, ru_text, eng_text, pt_br_text, uses, profession from translations where slug like  $1 || '%' and uses = true;`

	err = pgxscan.Select(ctx, tr.client, &translationsDao, q, slug)
	if err != nil {
		tr.logger.Error(fmt.Sprintf("Error while scanning row: %s op: %s", err, op))
		return translations, err
	}

	translations = make(map[string]appointment.TranslationEntity)
	for _, translation := range translationsDao {
		translations[translation.Slug] = translation.ToDomain()
	}
	return translations, nil
}

func (tr TranslationStorage) GetTranslationsBySlugKeyProfession(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error) {
	var translationsDao []dao.TranslationDao
	op := "sorkin_bot.internal.storage.read_repo.translations.GetTranslationsBySlugKeySlug"
	q := `select slug, ru_text, eng_text, pt_br_text, uses, profession from translations where slug like  $1 || '%' and uses = true;`

	err = pgxscan.Select(ctx, tr.client, &translationsDao, q, slug)
	if err != nil {
		tr.logger.Error(fmt.Sprintf("Error while scanning row: %s op: %s", err, op))
		return translations, err
	}

	translations = make(map[string]appointment.TranslationEntity)
	for _, translation := range translationsDao {
		translations[translation.Profession] = translation.ToDomain()
	}
	return translations, nil
}

func (tr TranslationStorage) GetTranslationsBySourceId(ctx context.Context, sourceId int) (translation appointment.TranslationEntity, err error) {
	var translationsDao dao.TranslationDao
	op := "sorkin_bot.internal.storage.read_repo.translations.GetTranslationsBySourceId"
	q := `select slug, ru_text, eng_text, pt_br_text, uses from translations where id_in_source_system = $1 and uses = true;`

	err = pgxscan.Get(ctx, tr.client, &translationsDao, q, sourceId)
	if err != nil {
		tr.logger.Error(fmt.Sprintf("Error while scanning row: %s op: %s, Source Id: %d", err, op, sourceId))
		return translation, err
	}
	translation = translationsDao.ToDomain()

	return translation, nil
}

func (tr TranslationStorage) GetManyTranslationsByIds(ctx context.Context, ids []int) (translations []appointment.TranslationEntity, err error) {
	var translationsDao []dao.TranslationDao
	op := "sorkin_bot.internal.storage.read_repo.translations.GetManyTranslationsByIds"
	if len(ids) == 0 {
		return nil, fmt.Errorf("ids slice is empty")
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	q := fmt.Sprintf(
		`select slug, ru_text, eng_text, pt_br_text, uses 
				from translations 
				where id_in_source_system in (%s) and uses = true;
		`, strings.Join(placeholders, ","),
	)

	err = pgxscan.Select(ctx, tr.client, &translationsDao, q, args...)
	if err != nil {
		tr.logger.Error(fmt.Sprintf("Error while scanning row: %s op: %s", err, op))
		return translations, err
	}

	for _, translation := range translationsDao {
		translations = append(translations, translation.ToDomain())
	}

	return translations, nil
}

package appointment

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
)

type ReadMessageRepo interface {
	GetTranslationsBySlugKeySlug(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error)
	GetTranslationsBySourceId(ctx context.Context, sourceId int) (translation appointment.TranslationEntity, err error)
	GetManyTranslationsByIds(ctx context.Context, ids []int) (translations []appointment.TranslationEntity, err error)
}

type ReadDraftAppointmentRepo interface {
	GetUserDraftAppointment(ctx context.Context, tgId int64) (draftAppointment appointment.DraftAppointment, err error)
}

type CreateDraftAppointmentUseCase interface {
	Execute(ctx context.Context, tgId int64) error
}

type UpdateDraftAppointmentDate interface {
	Execute(ctx context.Context, tgId int64, timeStart, timeEnd, date string) error
}

type UpdateDraftAppointmentStatus interface {
	Execute(ctx context.Context, tgId int64, appointmentId int) error
}

type UpdateDraftAppointmentIntField interface {
	Execute(ctx context.Context, tgId int64, fieldValue int, fieldName string) error
}

type CleanDraftAppointmentUseCase interface {
	Execute(ctx context.Context, tgId int64) error
}

type FastUpdateDraftAppointmentUseCase interface {
	Execute(ctx context.Context, tgId int64, draftAppointment appointment.DraftAppointment, created bool) error
}

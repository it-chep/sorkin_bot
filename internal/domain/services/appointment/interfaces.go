package appointment

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
)

type ReadMessageRepo interface {
	GetTranslationsBySlug(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error)
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
	Execute(ctx context.Context, tgId int64) error
}

type UpdateDraftAppointmentIntField interface {
	Execute(ctx context.Context, tgId int64, fieldValue int, fieldName string) error
}

type CleanDraftAppointmentUseCase interface {
	Execute(ctx context.Context, tgId int64) error
}

package update_user_patient_id

import (
	"context"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
)

type UpdateUserPatientIdUseCase struct {
	writeRepo WriteRepo
	logger    *slog.Logger
}

func NewUpdateUserPatientIdUseCase(writeRepo WriteRepo, logger *slog.Logger) UpdateUserPatientIdUseCase {
	return UpdateUserPatientIdUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc UpdateUserPatientIdUseCase) Execute(ctx context.Context, user entity.User, patientId int64) (err error) {
	return uc.writeRepo.UpdateUserPatientId(ctx, user, patientId)
}

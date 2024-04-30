package update_user_patient_id

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type WriteRepo interface {
	UpdateUserPatientId(ctx context.Context, user entity.User, patientId int) (err error)
}

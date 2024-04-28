package change_user_status

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type ChangeStatusUseCase struct {
}

func NewChangeStatusUseCase() ChangeStatusUseCase {
	return ChangeStatusUseCase{}
}

func (c ChangeStatusUseCase) Execute(ctx context.Context, user entity.User, state string) (err error) {
	//TODO implement me
	panic("implement me")
}

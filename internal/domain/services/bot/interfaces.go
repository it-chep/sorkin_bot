package bot

import (
	"context"
	userEntity "sorkin_bot/internal/domain/entity/user"
)

type messageService interface {
	GetMessage(ctx context.Context, user userEntity.User, name string) (messageText string, err error)
}

type readMessagesRepo interface {
	GetMessageByCondition()
}

type AdministratorHelpUseCase interface {
	Execute()
}

type CancelAppointmentUseCase interface {
	Execute()
}

package bot

import (
	"context"
	userEntity "sorkin_bot/internal/domain/entity/user"
)

type MessageService interface {
	GetMessage(ctx context.Context, user userEntity.User, name string) (messageText string, err error)
}

type ReadMessagesRepo interface {
	GetMessageByCondition()
}

type AdministratorHelpUseCase interface {
	Execute()
}

type CancelAppointmentUseCase interface {
	Execute()
}

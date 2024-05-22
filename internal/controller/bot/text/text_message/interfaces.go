package text_message

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

type appointmentService interface {
	// appointmeent interfaces in service and gateway
	CreateAppointment(ctx context.Context, user entity.User, callbackData string) (appointmentId *int)

	// user interfaces in service and gateway
	GetPatient(ctx context.Context, user entity.User) (result bool)
	CreatePatient(ctx context.Context, user entity.User) (result bool)

	GetDraftAppointment(ctx context.Context, tgId int64) (draftAppointment appointment.DraftAppointment, err error)
}

type userService interface {
	GetUser(ctx context.Context, tgId int64) (user entity.User, err error)
	UpdatePatientId(ctx context.Context, user entity.User, patientId *int) (err error)
	UpdateBirthDate(ctx context.Context, dto tg.TgUserDTO, birthDate string) (user entity.User, result bool, err error)
	UpdateFullName(ctx context.Context, dto tg.TgUserDTO, fullName string) (user entity.User, result bool, err error)
	UpdatePhone(ctx context.Context, dto tg.TgUserDTO, phone string) (user entity.User, result bool, err error)
}

type messageService interface {
	GetMessage(ctx context.Context, user entity.User, name string) (messageText string, err error)
}

type botGateway interface {
	SendConfirmAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, doctorId int)
	SendChangeLanguageMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
}

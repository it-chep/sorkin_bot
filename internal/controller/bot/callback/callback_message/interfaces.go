package callback

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

type AppointmentService interface {
	// appointmeent interfaces in service and gateway
	GetAppointments(ctx context.Context, user entity.User) (appointments []appointment.Appointment)
	GetAppointmentDetail(ctx context.Context, user entity.User, appointmentId int) (appointmentEntity appointment.Appointment)
	CreateAppointment(ctx context.Context, user entity.User, callbackData string) (appointmentId *int)
	ConfirmAppointment(ctx context.Context, appointmentId int) (result bool)
	CancelAppointment(ctx context.Context, user entity.User, appointmentId int) (result bool)
	RescheduleAppointment(ctx context.Context, user entity.User, appointmentId int, movedTo string) (result bool)

	// doctors interfaces in service and gateway
	GetDoctors(ctx context.Context, specialityId int) (doctors []appointment.Doctor)

	// schedules interfaces in service and gateway
	GetSchedules(ctx context.Context, doctorId int)
	GetFastAppointmentSchedules(ctx context.Context) (schedulesMap map[int][]appointment.Schedule)
}

type UserService interface {
	GetUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error)
	ChangeLanguage(ctx context.Context, dto tg.TgUserDTO, languageCode string) (user entity.User, err error)
	ChangeState(ctx context.Context, dto tg.TgUserDTO, state string) (user entity.User, err error)
}

type MessageService interface {
	GetMessage(ctx context.Context, user entity.User, name string) (messageText string, err error)
	SaveMessageLog(ctx context.Context, message tg.MessageDTO) (err error)
}

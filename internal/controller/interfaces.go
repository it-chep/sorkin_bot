package controller

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	GetSpecialities(ctx context.Context) (specialities []appointment.Speciality, err error)
	GetTranslatedSpecialities(ctx context.Context, user entity.User, specialities []appointment.Speciality, offset int) (translatedSpecialities map[int]string, unTranslatedSpecialities []string, err error)

	// schedules interfaces in service and gateway
	GetSchedules(ctx context.Context, doctorId int)
	GetFastAppointmentSchedules(ctx context.Context) (schedulesMap map[int][]appointment.Schedule)
	GetPatient(ctx context.Context, user entity.User) (result bool)
	CreatePatient(ctx context.Context, user entity.User) (result bool)
}

type UserService interface {
	UpdatePatientId(ctx context.Context, user entity.User, patientId *int) (err error)
	UpdateBirthDate(ctx context.Context, dto tg.TgUserDTO, birthDate string) (user entity.User, err error)
	UpdatePhone(ctx context.Context, dto tg.TgUserDTO, phone string) (user entity.User, err error)
	UpdateThirdName(ctx context.Context, dto tg.TgUserDTO, thirdName string) (user entity.User, err error)
	GetUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error)
	ChangeLanguage(ctx context.Context, dto tg.TgUserDTO, languageCode string) (user entity.User, err error)
	ChangeState(ctx context.Context, dto tg.TgUserDTO, state string) (user entity.User, err error)
	RegisterNewUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error)
}

type MessageService interface {
	GetMessage(ctx context.Context, user entity.User, name string) (messageText string, err error)
	SaveMessageLog(ctx context.Context, message tg.MessageDTO) (err error)
}

type BotService interface {
	ConfigureChangeLanguageMessage(ctx context.Context, user entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)
	ConfigureGetSpecialityMessage(
		ctx context.Context,
		userEntity entity.User,
		translatedSpecialities map[int]string,
		offset int,
	) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)
}

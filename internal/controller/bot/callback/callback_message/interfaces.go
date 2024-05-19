package callback

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

type AppointmentSpeciality interface {
	GetSpecialities(ctx context.Context) (specialities []appointment.Speciality, err error)
	GetTranslatedSpecialities(ctx context.Context, user entity.User, specialities []appointment.Speciality, offset int) (translatedSpecialities map[int]string, unTranslatedSpecialities []string, err error)
	TranslateSpecialityByID(ctx context.Context, user entity.User, specialityId int) (translatedSpeciality string, err error)
}

type draftAppointment interface {
	GetDraftAppointment(ctx context.Context, tgId int64) (draftAppointment appointment.DraftAppointment, err error)
	UpdateDraftAppointmentStatus(ctx context.Context, tgId int64)
	UpdateDraftAppointmentDate(ctx context.Context, tgId int64, timeStart, timeEnd, date string)
	UpdateDraftAppointmentIntField(ctx context.Context, tgId int64, intVal int, fieldName string)
	CleanDraftAppointment(ctx context.Context, tgId int64)
}

type appointmentService interface {
	// appointmeent interfaces in service and gateway
	GetAppointments(ctx context.Context, user entity.User) (appointments []appointment.Appointment)
	GetAppointmentDetail(ctx context.Context, user entity.User, appointmentId int) (appointmentEntity appointment.Appointment)
	CreateAppointment(ctx context.Context, user entity.User, callbackData string) (appointmentId *int)
	ConfirmAppointment(ctx context.Context, appointmentId int) (result bool)
	CancelAppointment(ctx context.Context, user entity.User, appointmentId int) (result bool)
	RescheduleAppointment(ctx context.Context, user entity.User, appointmentId int, movedTo string) (result bool)

	// doctors interfaces in service and gateway
	GetDoctors(ctx context.Context, tgId int64, offset int, specialityId *int) (doctorsMap map[int]string)

	// schedules interfaces in service and gateway
	GetSchedules(ctx context.Context, userEntity entity.User, doctorId *int) (schedulesMap []appointment.Schedule, err error)
	GetFastAppointmentSchedules(ctx context.Context) (schedulesMap map[int][]appointment.Schedule)

	AppointmentSpeciality
	draftAppointment
}

type userService interface {
	GetUser(ctx context.Context, tgId int64) (user entity.User, err error)
	ChangeLanguage(ctx context.Context, dto tg.TgUserDTO, languageCode string) (user entity.User, err error)
	ChangeState(ctx context.Context, tgId int64, state string) (user entity.User, err error)
}

type messageService interface {
	GetMessage(ctx context.Context, user entity.User, name string) (messageText string, err error)
}

type botService interface {
	ConfigureGetSpecialityMessage(
		ctx context.Context,
		userEntity entity.User,
		translatedSpecialities map[int]string,
		offset int,
	) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)

	ConfigureGetDoctorMessage(
		ctx context.Context,
		userEntity entity.User,
		doctors map[int]string,
		offset int,
	) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)

	ConfigureGetScheduleMessage(
		ctx context.Context,
		userEntity entity.User,
		schedules []appointment.Schedule,
		offset int,
	) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)

	ConfigureConfirmAppointmentMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)

	ConfigureGetPhoneMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.ReplyKeyboardMarkup)
}

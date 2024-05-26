package bot_gateway

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

type appointmentService interface {
	GetFastAppointmentSchedules(ctx context.Context) (randomDoctors map[int]appointment.Schedule)
	GetAppointments(ctx context.Context, user entity.User) (appointments []appointment.Appointment)
	GetTranslationString(langCode string, translationEntity appointment.TranslationEntity) (translatedSpeciality string)
	GetDoctorInfo(ctx context.Context, user entity.User, doctorId int) (doctorEntity appointment.Doctor, err error)
}

type messageService interface {
	SaveMessageLog(ctx context.Context, messageDTO tg.MessageDTO) (err error)
	GetMessage(ctx context.Context, user entity.User, name string) (messageText string, err error)
	GetTranslationsBySlugKeyProfession(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error)
}

type keyboardsInterface interface {
	ConfigureMainMenuMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)

	ConfigureChangeLanguageMessage(ctx context.Context, user entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)

	ConfigureFastAppointmentMessage(
		ctx context.Context,
		userEntity entity.User,
		schedulesMap map[int]appointment.Schedule,
	) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)

	ConfigureGetMyAppointmentsMessage(
		ctx context.Context,
		userEntity entity.User,
		appointments []appointment.Appointment,
		offset int,
	) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)

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

	ConfigureGetPhoneMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.ReplyKeyboardMarkup)
	ConfigureConfirmAppointmentMessage(ctx context.Context, userEntity entity.User, doctorId int) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)
	ConfigureGetScheduleMessage(
		ctx context.Context,
		userEntity entity.User,
		schedules []appointment.Schedule,
		offset int,
	) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)

	ConfigureAppointmentDetailMessage(
		ctx context.Context,
		userEntity entity.User,
		appointmentEntity appointment.Appointment,
	) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)

	ConfigureDoctorInfoMessage(ctx context.Context, userEntity entity.User, doctorId int) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)
}

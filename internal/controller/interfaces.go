package controller

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

type draftAppointment interface {
	GetDraftAppointment(ctx context.Context, tgId int64) (draftAppointment appointment.DraftAppointment, err error)
	UpdateDraftAppointmentStatus(ctx context.Context, tgId int64, appointmentId int)
	UpdateDraftAppointmentDate(ctx context.Context, tgId int64, timeStart, timeEnd, date string)
	UpdateDraftAppointmentIntField(ctx context.Context, tgId int64, intVal int, fieldName string)
	CreateDraftAppointment(ctx context.Context, tgId int64)
	CleanDraftAppointment(ctx context.Context, tgId int64)
	FastUpdateDraftAppointment(ctx context.Context, tgId int64, doctorId int, timeStart string, timeEnd string)
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

	GetSpecialities(ctx context.Context) (specialities []appointment.Speciality, err error)
	GetTranslatedSpecialities(ctx context.Context, user entity.User, specialities []appointment.Speciality, offset int) (translatedSpecialities map[int]string, unTranslatedSpecialities []string, err error)
	TranslateSpecialityByID(ctx context.Context, user entity.User, specialityId int) (translatedSpeciality string, err error)

	// schedules interfaces in service and gateway
	GetSchedules(ctx context.Context, userEntity entity.User, doctorId *int) (schedulesMap []appointment.Schedule, err error)
	GetFastAppointmentSchedules(ctx context.Context) (randomDoctors map[int]appointment.Schedule)
	GetPatient(ctx context.Context, user entity.User) (result bool)
	CreatePatient(ctx context.Context, user entity.User) (result bool)

	draftAppointment
}

type userService interface {
	UpdatePatientId(ctx context.Context, user entity.User, patientId *int) (err error)
	UpdateBirthDate(ctx context.Context, dto tg.TgUserDTO, birthDate string) (user entity.User, result bool, err error)
	UpdateFullName(ctx context.Context, dto tg.TgUserDTO, fullName string) (user entity.User, result bool, err error)
	UpdatePhone(ctx context.Context, dto tg.TgUserDTO, phone string) (user entity.User, result bool, err error)
	GetUser(ctx context.Context, tgId int64) (user entity.User, err error)
	ChangeLanguage(ctx context.Context, dto tg.TgUserDTO, languageCode string) (user entity.User, err error)
	ChangeState(ctx context.Context, tgId int64, state string) (user entity.User, err error)
	RegisterNewUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error)
}

type messageService interface {
	GetMessage(ctx context.Context, user entity.User, name string) (messageText string, err error)
	SaveMessageLog(ctx context.Context, messageDTO tg.MessageDTO) (err error)
}

type botGateway interface {
	SendGetDoctorsMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, doctors map[int]string, offset int)
	SendChooseSpecialityMessage(
		ctx context.Context,
		idToDelete int,
		translatedSpecialities map[int]string,
		user entity.User,
		messageDTO tg.MessageDTO,
	)
	SendStartMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendChangeLanguageMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendGetPhoneMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendMyAppointmentsMessage(ctx context.Context, user entity.User, appointments []appointment.Appointment, messageDTO tg.MessageDTO)
	SendConfirmAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, doctorId int)
	SendFastAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendDetailAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, appointmentId int)
	SendSchedulesMessage(ctx context.Context, userEntity entity.User, messageDTO tg.MessageDTO, schedules []appointment.Schedule, offset int)
	SendSpecialityMessage(ctx context.Context, userEntity entity.User, messageDTO tg.MessageDTO, specialities map[int]string, offset int)
}

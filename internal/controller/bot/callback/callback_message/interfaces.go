package callback

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"time"
)

type AppointmentSpeciality interface {
	GetSpecialities(ctx context.Context) (specialities []appointment.Speciality, err error)
	GetTranslatedSpecialities(ctx context.Context, user entity.User, specialities []appointment.Speciality, offset int) (translatedSpecialities map[int]string, unTranslatedSpecialities []string, err error)
	TranslateSpecialityByID(ctx context.Context, user entity.User, specialityId int) (translatedSpeciality string, err error)
}

type draftAppointment interface {
	CreateDraftAppointment(ctx context.Context, tgId int64)
	GetDraftAppointment(ctx context.Context, tgId int64) (draftAppointment appointment.DraftAppointment, err error)
	UpdateDraftAppointmentStatus(ctx context.Context, tgId int64, appointmentId int)
	UpdateDraftAppointmentDate(ctx context.Context, tgId int64, timeStart, timeEnd, date string)
	UpdateDraftAppointmentIntField(ctx context.Context, tgId int64, intVal int, fieldName string)
	UpdateDraftAppointmentDoctorName(ctx context.Context, tgId int64, doctorId int)
	UpdateDraftAppointmentType(ctx context.Context, tgId int64, appointmentType appointment.AppointmentType)
	CleanDraftAppointment(ctx context.Context, tgId int64)
	FastUpdateDraftAppointment(ctx context.Context, tgId int64, specialityId, doctorId int, timeStart, timeEnd string)
	GetDraftAppointmentByAppointmentId(ctx context.Context, appointmentId int) (draftAppointment appointment.DraftAppointment, err error)
}

type appointmentService interface {
	// appointmeent interfaces in service and gateway
	GetAppointments(ctx context.Context, user entity.User) (appointments []appointment.Appointment)
	GetAppointmentDetail(ctx context.Context, user entity.User, appointmentId int) (appointmentEntity appointment.Appointment)
	CreateAppointment(ctx context.Context, user entity.User, draftAppointment appointment.DraftAppointment, callbackData string) (appointmentId *int)
	ConfirmAppointment(ctx context.Context, appointmentId int) (result bool)
	CancelAppointment(ctx context.Context, user entity.User, appointmentId int) (result bool)
	RescheduleAppointment(ctx context.Context, user entity.User, appointmentId int, movedTo string) (result bool)

	// doctors interfaces in service and gateway
	GetDoctorsBySpecialityId(ctx context.Context, tgId int64, offset int, specialityId *int) (doctorsMap map[int]string)
	GetDoctors(ctx context.Context, tgId int64, offset int) (doctorsMap map[int]string)

	// schedules interfaces in service and gateway
	GetSchedulesByDoctorId(ctx context.Context, userEntity entity.User, dayStart time.Time, doctorId *int) (schedulesMap []appointment.Schedule, err error)
	GetSchedulesToHomeVisit(ctx context.Context, userEntity entity.User, dayStart time.Time) (schedulesMap []appointment.Schedule, err error)
	GetFastAppointmentSchedules(ctx context.Context) (randomDoctors map[int]appointment.Schedule)
	GetSchedulePeriodsByDoctorId(ctx context.Context, doctorId int, dayStart time.Time) (schedulePeriodsMap map[time.Time]bool, err error)

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

type botGateway interface {
	SendStartMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendChangeLanguageMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendMyAppointmentsMessage(ctx context.Context, user entity.User, appointments []appointment.Appointment, messageDTO tg.MessageDTO, offset int)
	SendFastAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendGetDoctorsMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, doctors map[int]string, offset int)
	SendConfirmAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, doctorId int)
	SendGetPhoneMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendGetHomeAddressMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendSchedulesMessage(ctx context.Context, userEntity entity.User, messageDTO tg.MessageDTO, schedules []appointment.Schedule, offset int)
	SendSpecialityMessage(ctx context.Context, userEntity entity.User, messageDTO tg.MessageDTO, specialities map[int]string, offset int)
	SendDetailAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, appointmentEntity appointment.Appointment)
	SendEmptyAppointments(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendWaitMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, waitMessage string) int
	SendError(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendDoctorInfoMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, idToDelete, doctorId int)
	SendCalendarMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, year int, month time.Month, schedulesMap map[time.Time]bool)
	SendForbiddenAction(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendChooseAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendDoctorOrReasonMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendHomeDoctorSpecificationMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendEmptySchedulesHomeVisit(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendEmptySchedulePeriods(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
}

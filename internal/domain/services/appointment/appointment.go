package appointment

import (
	"context"
	"fmt"
	"log/slog"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/services/adapter"
	"sorkin_bot/internal/domain/services/user"
	"strconv"
	"strings"
)

type AppointmentService struct {
	misAdapter                        *adapter.AppointmentServiceAdapter
	userService                       user.UserService
	logger                            *slog.Logger
	readMessageRepo                   ReadMessageRepo
	readDraftAppointmentRepo          ReadDraftAppointmentRepo
	createDraftAppointmentUseCase     CreateDraftAppointmentUseCase
	updateDraftAppointmentDate        UpdateDraftAppointmentDate
	updateDraftAppointmentStatus      UpdateDraftAppointmentStatus
	updateDraftAppointmentIntField    UpdateDraftAppointmentIntField
	cleanDraftAppointmentUseCase      CleanDraftAppointmentUseCase
	fastUpdateDraftAppointmentUseCase FastUpdateDraftAppointmentUseCase
}

func NewAppointmentService(
	misAdapter *adapter.AppointmentServiceAdapter,
	readMessageRepo ReadMessageRepo,
	readDraftAppointmentRepo ReadDraftAppointmentRepo,
	logger *slog.Logger,
	userService user.UserService,
	createDraftAppointmentUseCase CreateDraftAppointmentUseCase,
	updateDraftAppointmentDate UpdateDraftAppointmentDate,
	updateDraftAppointmentStatus UpdateDraftAppointmentStatus,
	updateDraftAppointmentIntField UpdateDraftAppointmentIntField,
	cleanDraftAppointmentUseCase CleanDraftAppointmentUseCase,
	fastUpdateDraftAppointmentUseCase FastUpdateDraftAppointmentUseCase,
) AppointmentService {
	return AppointmentService{
		misAdapter:                        misAdapter,
		userService:                       userService,
		readMessageRepo:                   readMessageRepo,
		readDraftAppointmentRepo:          readDraftAppointmentRepo,
		logger:                            logger,
		createDraftAppointmentUseCase:     createDraftAppointmentUseCase,
		updateDraftAppointmentDate:        updateDraftAppointmentDate,
		updateDraftAppointmentStatus:      updateDraftAppointmentStatus,
		updateDraftAppointmentIntField:    updateDraftAppointmentIntField,
		cleanDraftAppointmentUseCase:      cleanDraftAppointmentUseCase,
		fastUpdateDraftAppointmentUseCase: fastUpdateDraftAppointmentUseCase,
	}
}

func (as *AppointmentService) GetAppointments(ctx context.Context, user entity.User) (appointments []appointment.Appointment) {
	_ = "sorkin_bot.internal.domain.services.appointment.appointment.GetAppointments"
	if user.GetPatientId() == nil {
		return
	}

	return as.misAdapter.MyAppointments(ctx, user)
}

func (as *AppointmentService) GetAppointmentDetail(ctx context.Context, user entity.User, appointmentId int) (appointmentEntity appointment.Appointment) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.GetAppointmentDetail"

	if user.GetPatientId() == nil {
		return
	}

	appointmentEntity, err := as.misAdapter.DetailAppointment(ctx, user, appointmentId)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return appointment.Appointment{}
	}

	return appointmentEntity
}

func (as *AppointmentService) CreateAppointment(ctx context.Context, user entity.User, callbackData string) (appointmentId *int) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.CreateAppointment"
	if user.GetPatientId() == nil {
		return
	}
	// example: callbackData = doctorId_8__timeStart_'11.05.2004 12:00'__timeEnd_'11.05.2004 12:30'
	elements := strings.Split(callbackData, "__")

	doctorId, _ := strconv.Atoi(strings.Split(elements[0], "_")[1])
	timeStart, timeEnd := as.convertToValidDate(elements)

	appointmentId, err := as.misAdapter.CreateAppointment(ctx, user, doctorId, timeStart, timeEnd)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return nil
	}

	return appointmentId
}

func (as *AppointmentService) convertToValidDate(elements []string) (timeStartValid, timeEndValid string) {
	timeStartDirt := strings.Split(elements[1], "_")[1]
	timeEndDirt := strings.Split(elements[2], "_")[1]

	dateTimeStart := strings.Split(timeStartDirt, " ")
	dateTimeEnd := strings.Split(timeEndDirt, " ")

	date := strings.Split(dateTimeStart[0], "-")
	day, _ := strconv.Atoi(date[2])
	month, _ := strconv.Atoi(date[1])
	year, _ := strconv.Atoi(date[0])

	timeStart := strings.Split(dateTimeStart[1], ":")
	hourStart, _ := strconv.Atoi(timeStart[0])
	minuteStart, _ := strconv.Atoi(timeStart[1])

	timeEnd := strings.Split(dateTimeEnd[1], ":")
	hourEnd, _ := strconv.Atoi(timeEnd[0])
	minuteEnd, _ := strconv.Atoi(timeEnd[1])

	timeStartValid = fmt.Sprintf("%02d.%02d.%04d %02d:%02d", day, month, year, hourStart, minuteStart)
	timeEndValid = fmt.Sprintf("%02d.%02d.%04d %02d:%02d", day, month, year, hourEnd, minuteEnd)

	return timeStartValid, timeEndValid
}

func (as *AppointmentService) ConfirmAppointment(ctx context.Context, appointmentId int) (result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.ConfirmAppointment"

	result, err := as.misAdapter.ConfirmAppointment(ctx, appointmentId)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return false
	}
	return true
}

func (as *AppointmentService) CancelAppointment(ctx context.Context, user entity.User, appointmentId int) (result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.CancelAppointment"
	result, err := as.misAdapter.CancelAppointment(ctx, user, appointmentId)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return false
	}
	return true
}

func (as *AppointmentService) RescheduleAppointment(ctx context.Context, user entity.User, appointmentId int, movedTo string) (result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.RescheduleAppointment"

	err := as.misAdapter.RescheduleAppointment(ctx, user, movedTo, appointmentId)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return false
	}
	return true
}

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
	misAdapter                     *adapter.AppointmentServiceAdapter
	userService                    user.UserService
	logger                         *slog.Logger
	readMessageRepo                ReadMessageRepo
	readDraftAppointmentRepo       ReadDraftAppointmentRepo
	createDraftAppointmentUseCase  CreateDraftAppointmentUseCase
	updateDraftAppointmentDate     UpdateDraftAppointmentDate
	updateDraftAppointmentStatus   UpdateDraftAppointmentStatus
	updateDraftAppointmentIntField UpdateDraftAppointmentIntField
	cleanDraftAppointmentUseCase   CleanDraftAppointmentUseCase
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
) AppointmentService {
	return AppointmentService{
		misAdapter:                     misAdapter,
		userService:                    userService,
		readMessageRepo:                readMessageRepo,
		readDraftAppointmentRepo:       readDraftAppointmentRepo,
		logger:                         logger,
		createDraftAppointmentUseCase:  createDraftAppointmentUseCase,
		updateDraftAppointmentDate:     updateDraftAppointmentDate,
		updateDraftAppointmentStatus:   updateDraftAppointmentStatus,
		updateDraftAppointmentIntField: updateDraftAppointmentIntField,
		cleanDraftAppointmentUseCase:   cleanDraftAppointmentUseCase,
	}
}

func (as *AppointmentService) GetAppointments(ctx context.Context, user entity.User) (appointments []appointment.Appointment) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.GetAppointments"
	var err error
	if user.GetPatientId() == nil {
		return
	}

	appointments = as.misAdapter.MyAppointments(ctx, user)

	for _, appointmentEntity := range appointments {
		as.logger.Info(fmt.Sprintf("%d %s %s", appointmentEntity.GetAppointmentId(), appointmentEntity.GetTimeStart(), op))
	}

	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return
	}

	return appointments
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
	timeStart := strings.Split(elements[1], "_")[1]
	timeEnd := strings.Split(elements[2], "_")[1]

	appointmentId, err := as.misAdapter.CreateAppointment(ctx, user, doctorId, timeStart, timeEnd)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return nil
	}

	return appointmentId
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

package appointment

import (
	"context"
	"fmt"
	"log/slog"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/services/user"
	"strconv"
	"strings"
)

type AppointmentService struct {
	mis         Appointment
	userService user.UserService
	readRepo    ReadRepo
	logger      *slog.Logger
}

func NewAppointmentService(mis Appointment, readRepo ReadRepo, logger *slog.Logger, userService user.UserService) AppointmentService {
	return AppointmentService{
		mis:         mis,
		userService: userService,
		readRepo:    readRepo,
		logger:      logger,
	}
}

func (as *AppointmentService) GetAppointments(ctx context.Context, user entity.User) (appointments []appointment.Appointment) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.GetAppointments"

	if user.GetPatientId() == 0 {
		return
	}

	err, appointments := as.mis.MyAppointments(ctx, user)

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

	if user.GetPatientId() == 0 {
		return
	}

	err, appointmentEntity := as.mis.DetailAppointment(ctx, user, appointmentId)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return appointment.Appointment{}
	}

	return appointmentEntity
}

func (as *AppointmentService) CreateAppointment(ctx context.Context, user entity.User, callbackData string) (appointmentId int) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.CreateAppointment"
	if user.GetPatientId() == 0 {
		return
	}
	// example: callbackData = doctorId_8__timeStart_'11.05.2004 12:00'__timeEnd_'11.05.2004 12:30'
	elements := strings.Split(callbackData, "__")

	doctorId, _ := strconv.Atoi(strings.Split(elements[0], "_")[1])
	timeStart := strings.Split(elements[1], "_")[1]
	timeEnd := strings.Split(elements[2], "_")[1]

	err, appointmentId := as.mis.CreateAppointment(ctx, user, doctorId, timeStart, timeEnd)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return -1
	}

	return appointmentId
}

func (as *AppointmentService) ConfirmAppointment(ctx context.Context, appointmentId int) (result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.ConfirmAppointment"

	err, result := as.mis.ConfirmAppointment(ctx, appointmentId)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return false
	}
	return true
}

func (as *AppointmentService) CancelAppointment(ctx context.Context, appointmentId int) (result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.CancelAppointment"

	err, result := as.mis.CancelAppointment(ctx, "", appointmentId)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return false
	}
	return true
}

func (as *AppointmentService) RescheduleAppointment(ctx context.Context, appointmentId int, movedTo string) (result bool) {
	op := "sorkin_bot.internal.domain.services.appointment.appointment.RescheduleAppointment"

	err, result := as.mis.CancelAppointment(ctx, movedTo, appointmentId)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return false
	}
	return true
}

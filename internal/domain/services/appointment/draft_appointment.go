package appointment

import (
	"context"
	"fmt"
	"reflect"
	"sorkin_bot/internal/domain/entity/appointment"
)

func (as *AppointmentService) GetDraftAppointment(ctx context.Context, tgId int64) (draftAppointment appointment.DraftAppointment, err error) {
	draftAppointment, err = as.readDraftAppointmentRepo.GetUserDraftAppointment(ctx, tgId)
	if err != nil {
		return appointment.NewDraftAppointment(
			nil, nil, nil, nil, nil, nil, nil, nil,
		), err
	}
	return draftAppointment, nil
}

func (as *AppointmentService) CreateDraftAppointment(ctx context.Context, tgId int64) {
	draftAppointment, err := as.GetDraftAppointment(ctx, tgId)

	if err != nil {
		return
	}

	if !reflect.ValueOf(draftAppointment.GetTgId()).IsNil() {
		return
	}

	err = as.createDraftAppointmentUseCase.Execute(ctx, tgId)
	if err != nil {
		return
	}
}

func (as *AppointmentService) UpdateDraftAppointmentStatus(ctx context.Context, tgId int64, appointmentId int) {
	err := as.updateDraftAppointmentStatus.Execute(ctx, tgId, appointmentId)
	if err != nil {
		return
	}
}

func (as *AppointmentService) UpdateDraftAppointmentDate(ctx context.Context, tgId int64, timeStart, timeEnd, date string) {
	err := as.updateDraftAppointmentDate.Execute(ctx, tgId, timeStart, timeEnd, date)
	if err != nil {
		return
	}
}

func (as *AppointmentService) UpdateDraftAppointmentIntField(ctx context.Context, tgId int64, intVal int, fieldName string) {
	err := as.updateDraftAppointmentIntField.Execute(ctx, tgId, intVal, fieldName)
	if err != nil {
		return
	}
}

func (as *AppointmentService) UpdateDraftAppointmentDoctorName(ctx context.Context, tgId int64, doctorId int) {
	doctor := as.misAdapter.GetDoctorInfo(ctx, doctorId)
	err := as.updateDraftAppointmentStrField.Execute(ctx, tgId, doctor.GetName(), "doctor_name")
	if err != nil {
		return
	}
}

func (as *AppointmentService) CleanDraftAppointment(ctx context.Context, tgId int64) {
	err := as.cleanDraftAppointmentUseCase.Execute(ctx, tgId)
	if err != nil {
		return
	}
}

func (as *AppointmentService) FastUpdateDraftAppointment(ctx context.Context, tgId int64, specialityId, doctorId int, timeStart, timeEnd string) {
	var created = true
	draftAppointment := appointment.NewDraftAppointment(&specialityId, &doctorId, &tgId, nil, &timeStart, &timeEnd, nil, nil)
	oldDraftAppointment, err := as.readDraftAppointmentRepo.GetUserDraftAppointment(ctx, tgId)
	if err != nil {
		return
	}
	if oldDraftAppointment.GetTgId() == nil {
		created = false
	}
	err = as.fastUpdateDraftAppointmentUseCase.Execute(ctx, tgId, draftAppointment, created)
	if err != nil {
		as.logger.Error(fmt.Sprintf("fast update draft appointment failed: %s", err))
		return
	}
}

func (as *AppointmentService) GetDraftAppointmentByAppointmentId(ctx context.Context, appointmentId int) (draftAppointment appointment.DraftAppointment, err error) {
	draftAppointment, err = as.readDraftAppointmentRepo.GetDraftAppointmentByAppointmentId(ctx, appointmentId)
	if err != nil {
		return appointment.NewDraftAppointment(
			nil, nil, nil, nil, nil, nil, nil, nil,
		), err
	}
	return draftAppointment, nil
}

func (as *AppointmentService) UpdateDraftAppointmentType(ctx context.Context, tgId int64, appointmentType appointment.AppointmentType) {
	err := as.updateDraftAppointmentStrField.Execute(ctx, tgId, string(appointmentType), "type")
	if err != nil {
		return
	}
}

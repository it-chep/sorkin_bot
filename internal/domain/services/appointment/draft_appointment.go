package appointment

import (
	"context"
	"reflect"
	"sorkin_bot/internal/domain/entity/appointment"
)

func (as *AppointmentService) GetDraftAppointment(ctx context.Context, tgId int64) (draftAppointment appointment.DraftAppointment, err error) {
	draftAppointment, err = as.readDraftAppointmentRepo.GetUserDraftAppointment(ctx, tgId)
	if err != nil {
		return appointment.NewDraftAppointment(
			nil, nil, nil, nil, nil, nil,
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

func (as *AppointmentService) CleanDraftAppointment(ctx context.Context, tgId int64) {
	err := as.cleanDraftAppointmentUseCase.Execute(ctx, tgId)
	if err != nil {
		return
	}
}

func (as *AppointmentService) FastUpdateDraftAppointment(ctx context.Context, tgId int64, doctorId int, timeStart, timeEnd string) {
	draftAppointment := appointment.NewDraftAppointment(nil, &doctorId, &tgId, &timeStart, &timeEnd, nil)
	err := as.fastUpdateDraftAppointmentUseCase.Execute(ctx, tgId, draftAppointment)
	if err != nil {
		return
	}
}

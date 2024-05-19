package appointment

import (
	"context"
	"reflect"
	"sorkin_bot/internal/domain/entity/appointment"
)

func (as *AppointmentService) GetDraftAppointment(ctx context.Context, tgId int64) (draftAppointment appointment.DraftAppointment, err error) {
	draftAppointment, err = as.readDraftAppointmentRepo.GetUserDraftAppointment(ctx, tgId)
	if err != nil {
		return appointment.DraftAppointment{}, err
	}
	return draftAppointment, nil
}

func (as *AppointmentService) CreateDraftAppointment(ctx context.Context, tgId int64) {
	draftAppointment, err := as.GetDraftAppointment(ctx, tgId)
	if err != nil {
		return
	}
	if !reflect.ValueOf(draftAppointment.GetTgId()).IsZero() {
		return
	}

	err = as.createDraftAppointmentUseCase.Execute(ctx, tgId)
	if err != nil {
		return
	}
}

func (as *AppointmentService) UpdateDraftAppointmentStatus(ctx context.Context, tgId int64) {
	err := as.updateDraftAppointmentStatus.Execute(ctx, tgId)
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

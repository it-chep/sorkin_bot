package write_repo

import (
	"context"
	"fmt"
	"log/slog"
	"sorkin_bot/internal/domain/entity/appointment"
	"sorkin_bot/pkg/client/postgres"
)

type AppointmentStorage struct {
	client postgres.Client
	logger *slog.Logger
}

func NewAppointmentStorage(client postgres.Client, logger *slog.Logger) AppointmentStorage {
	return AppointmentStorage{
		client: client,
		logger: logger,
	}
}

func (rs AppointmentStorage) CreateEmptyDraftAppointment(ctx context.Context, tgId int64) (err error) {
	op := "internal/storage/read_repo/appointment/CreateEmptyDraftAppointment"
	q := `
		insert into appointment (tg_id, draft) values ($1, true) returning id;
	`
	var draftAppointmentId *int
	err = rs.client.QueryRow(ctx, q, tgId).Scan(&draftAppointmentId)
	if err != nil {
		rs.logger.Error(fmt.Sprintf("Error while scanning row: %s, op: %s", err, op))
		return err
	}

	return nil
}

func (rs AppointmentStorage) UpdateDateDraftAppointment(
	ctx context.Context, tgId int64, timeStart, timeEnd, date string,
) (err error) {
	op := "internal/storage/read_repo/appointment/UpdateDateDraftAppointment"
	q := `
		update appointment set time_start = $1, time_end = $2, date = $3 where tg_id = $4 and draft = true;
	`

	_, err = rs.client.Exec(ctx, q, timeStart, timeEnd, date, tgId)
	if err != nil {
		rs.logger.Error(fmt.Sprintf("Error while executing row: %s, op: %s", err, op))
		return err
	}

	return nil
}

func (rs AppointmentStorage) UpdateIntDraftAppointment(ctx context.Context, tgId int64, intValue int, intField string) (err error) {
	op := "internal/storage/read_repo/appointment/UpdateIntDraftAppointment"
	q := fmt.Sprintf(`update appointment set %s = $1 where tg_id = $2 and draft = true;`, intField)

	_, err = rs.client.Exec(ctx, q, intValue, tgId)
	if err != nil {
		rs.logger.Error(fmt.Sprintf("Error while executing row: %s, op: %s", err, op))
		return err
	}

	return nil
}

func (rs AppointmentStorage) UpdateStatusDraftAppointment(ctx context.Context, tgId int64, appointmentId int) (err error) {
	op := "internal/storage/read_repo/appointment/UpdateStatusDraftAppointment"
	q := `
		update appointment set appointment_id = $1, draft = false where tg_id = $2 and draft = true;
	`

	_, err = rs.client.Exec(ctx, q, appointmentId, tgId)
	if err != nil {
		rs.logger.Error(fmt.Sprintf("Error while executing row: %s, op: %s", err, op))
		return err
	}

	return nil
}

func (rs AppointmentStorage) FastUpdateDraftAppointment(
	ctx context.Context, tgId int64,
	draftAppointment appointment.DraftAppointment,
) (err error) {
	op := "internal/storage/read_repo/appointment/FastUpdateDraftAppointment"
	q := `
		insert into appointment (tg_id, doctor_id, time_start, time_end, draft)
		values ($1, $2, $3, $4, true);
	`

	var doctorId int
	if draftAppointment.GetDoctorId() != nil {
		doctorId = *draftAppointment.GetDoctorId()
	}

	var timeStart, timeEnd string
	if draftAppointment.GetTimeStart() != nil {
		timeStart = *draftAppointment.GetTimeStart()
	}
	if draftAppointment.GetTimeEnd() != nil {
		timeEnd = *draftAppointment.GetTimeEnd()
	}

	_, err = rs.client.Exec(
		ctx, q,
		tgId,
		doctorId,
		timeStart,
		timeEnd,
	)

	if err != nil {
		rs.logger.Error(fmt.Sprintf("Error while executing row: %s, op: %s", err, op))
		return err
	}

	return nil
}

func (rs AppointmentStorage) CleanDraftAppointment(ctx context.Context, tgId int64) (err error) {
	op := "internal/storage/read_repo/appointment/ClearDraftAppointment"
	q := `
		update appointment 
		set 
		    speciality_id=null, 
		    doctor_id=null, 
		    time_start=null, 
		    time_end=null, 
		    date=null 
		where tg_id = $1 and draft = true;
	`
	_, err = rs.client.Exec(ctx, q, tgId)
	if err != nil {
		rs.logger.Error(fmt.Sprintf("Error while executing row: %s, op: %s", err, op))
		return err
	}

	return nil
}

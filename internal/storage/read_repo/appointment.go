package read_repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"sorkin_bot/internal/domain/entity/appointment"
	"sorkin_bot/internal/storage/dao"
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

func (rs AppointmentStorage) GetUserDraftAppointment(ctx context.Context, tgId int64) (draftAppointment appointment.DraftAppointment, err error) {
	op := "internal/storage/read_repo/appointment/GetUserDraftAppointment"
	q := `
		select id, tg_id, speciality_id, doctor_id, date, time_start, time_end, draft 
		from appointment
		where tg_id = $1 and draft = true;
	`
	var appointmentDAO dao.AppointmentDAO
	err = pgxscan.Get(ctx, rs.client, &appointmentDAO, q, tgId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appointment.DraftAppointment{}, nil
		}
		rs.logger.Error(fmt.Sprintf("Error while scanning row: %s, op: %s", err, op))
		return appointment.DraftAppointment{}, err
	}

	return appointmentDAO.ToDomain(), nil
}

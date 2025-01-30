package notify_appointment

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron"
	"log/slog"
	"os"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/pkg/client/telegram"
	"strconv"
	"time"
)

var schedule, _ = cron.ParseStandard("*/15 * * * *")

type Task struct {
	appointmentService  appointmentService
	notificationService notificationService
	logger              *slog.Logger
	bot                 telegram.Bot
}

func NewTask(
	appointmentService appointmentService,
	notificationService notificationService,
	logger *slog.Logger,
	bot telegram.Bot,
) Task {
	return Task{
		appointmentService:  appointmentService,
		notificationService: notificationService,
		logger:              logger,
		bot:                 bot,
	}
}

func (t Task) Process(ctx context.Context) error {
	adminId, err := strconv.Atoi(os.Getenv("ADMIN_ID"))
	if err != nil {
		panic("adminId not found")
	}

	appointments, err := t.appointmentService.GetAppointmentsForNotifying(ctx)
	if err != nil {
		return err
	}

	t.logger.Info("Processing appointments", time.Now())

	for _, appointmentEntity := range appointments {
		t.logger.Info(fmt.Sprintf(
			"fake send to %s, appointment id %d",
			appointmentEntity.GetPatientPhone(),
			appointmentEntity.GetAppointmentId()),
		)

		err = t.notificationService.NotifySoonAppointment(ctx, appointmentEntity)
		if err != nil {
			t.notifyErrorTelegramAdmin(int64(adminId), err)
			return err
		}

		// fake break because spam
		break
	}

	return nil
}

func (t Task) notifyErrorTelegramAdmin(adminId int64, err error) {
	messageDTO := tg.MessageDTO{
		Chat: &tg.Chat{
			ID: adminId,
		},
	}

	msg := tgbotapi.NewMessage(
		adminId,
		fmt.Sprintf("Ошибка при отправке сообщения пользвателю по СМС о скором визите в Клинику %s", err),
	)

	t.bot.SendMessage(msg, messageDTO)
}

func (t Task) NextSchedule(now time.Time) time.Time {
	return schedule.Next(now)
}

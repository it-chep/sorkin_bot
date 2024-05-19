package tasks

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/pkg/client/telegram"
	"strconv"
)

type GetTranslatedSpecialityTask struct {
	appointmentService appointmentService
	userService        userService
	logger             *slog.Logger
	bot                telegram.Bot
}

func NewGetTranslatedSpecialityTask(appointmentService appointmentService, userService userService, logger *slog.Logger, bot telegram.Bot) GetTranslatedSpecialityTask {
	return GetTranslatedSpecialityTask{
		appointmentService: appointmentService,
		userService:        userService,
		logger:             logger,
		bot:                bot,
	}
}

func (task GetTranslatedSpecialityTask) Process(ctx context.Context) error {
	adminId, err := strconv.Atoi(os.Getenv("ADMIN_ID"))
	return nil
	if err != nil {
		panic("adminId not found")
	}
	dto := tg.TgUserDTO{TgID: int64(adminId)}
	chatId := tg.Chat{ID: int64(adminId)}
	messageDTO := tg.MessageDTO{
		Chat: &chatId,
	}

	getUser, err := task.userService.GetUser(ctx, dto.TgID)
	if err != nil {
		msg := tgbotapi.NewMessage(int64(adminId), "error while getting admin in GetTranslatedSpecialityTask")
		task.bot.SendMessage(msg, messageDTO)
		return err
	}
	specialities, err := task.appointmentService.GetSpecialities(ctx)
	if err != nil {
		msg := tgbotapi.NewMessage(int64(adminId), "error while getting speciality in GetTranslatedSpecialityTask")
		task.bot.SendMessage(msg, messageDTO)
		return err
	}
	_, unTranslatedSpecialities, err := task.appointmentService.GetTranslatedSpecialities(ctx, getUser, specialities, 0)
	for _, unTranslatedSpeciality := range unTranslatedSpecialities {
		msg := tgbotapi.NewMessage(int64(adminId), fmt.Sprintf("untranslated speciality %s !", unTranslatedSpeciality))
		task.bot.SendMessage(msg, messageDTO)
	}
	return nil
}

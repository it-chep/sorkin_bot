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

type CheckAdministrationHelpTask struct {
	logger         *slog.Logger
	bot            telegram.Bot
	messageService messageService
	userService    userService
}

func NewCheckAdministrationHelpTask(logger *slog.Logger, bot telegram.Bot, messageService messageService, userService userService) CheckAdministrationHelpTask {
	return CheckAdministrationHelpTask{
		logger:         logger,
		bot:            bot,
		messageService: messageService,
		userService:    userService,
	}
}

func (task CheckAdministrationHelpTask) Process(ctx context.Context) error {
	adminId, err := strconv.Atoi(os.Getenv("ADMIN_ID"))
	if err != nil {
		panic("adminId not found")
	}
	minutes, err := strconv.Atoi(os.Getenv("DEFAULT_CHECK_SUPPORT"))
	if err != nil {
		panic("DEFAULT_CHECK_SUPPORT not found")
	}

	chatId := tg.Chat{ID: int64(adminId)}
	messageDTO := tg.MessageDTO{
		Chat: &chatId,
	}
	logs, err := task.messageService.GetSupportLogs(ctx, minutes)
	if err != nil {
		return err
	}

	for _, log := range logs {
		userEntity, err := task.userService.GetUser(ctx, log.GetUserTgId())
		if err != nil {
			continue
		}
		msg := tgbotapi.NewMessage(int64(adminId),
			fmt.Sprintf(
				"tech_support call from %s %s tg_id: %d", userEntity.GetFirstName(), userEntity.GetLastName(), userEntity.GetTgId(),
			),
		)
		task.bot.SendMessage(msg, messageDTO)
	}
	return nil
}

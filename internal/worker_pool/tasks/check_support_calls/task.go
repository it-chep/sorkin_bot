package check_support_calls

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/pkg/client/telegram"
	"strconv"
	"time"
)

type Task struct {
	logger         *slog.Logger
	bot            telegram.Bot
	messageService messageService
	userService    userService
}

func NewTask(
	logger *slog.Logger,
	bot telegram.Bot,
	messageService messageService,
	userService userService,
) Task {
	return Task{
		logger:         logger,
		bot:            bot,
		messageService: messageService,
		userService:    userService,
	}
}

func (t Task) Process(ctx context.Context) error {
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

	logs, err := t.messageService.GetSupportLogs(ctx, minutes)
	if err != nil {
		return err
	}

	for _, log := range logs {
		userEntity, err := t.userService.GetUser(ctx, log.GetUserTgId())
		if err != nil {
			continue
		}

		msg := tgbotapi.NewMessage(int64(adminId),
			fmt.Sprintf(
				"tech_support call from %s %s tg_id: %d", userEntity.GetFirstName(), userEntity.GetLastName(), userEntity.GetTgId(),
			),
		)

		t.bot.SendMessage(msg, messageDTO)
	}
	return nil
}

func (t Task) NextSchedule(now time.Time) time.Time {
	return now.Add(5 * time.Minute)
}
